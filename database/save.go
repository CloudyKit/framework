package database

import (
	"github.com/CloudyKit/framework/database/change"
	"github.com/CloudyKit/framework/database/model"
	"github.com/CloudyKit/framework/database/scheme"
	"github.com/CloudyKit/framework/validation"

	"errors"
	"fmt"
	"reflect"
)

func (db *DB) updateModelFields(m, mp Model, mpField string, mRef reflect.Value, v *validation.Validator) (*validation.Validator, error) {
	var err error

	mDat, mpDat := model.GetModelData(m), model.GetModelData(mp)

	fields := mDat.Scheme.Fields()
	operations := make([]change.Set, len(fields))
	var i int

	for _, field := range fields {

		mField := mRef.FieldByName(field.Name)
		if !mField.IsValid() {
			if field.Required {
				return v, fmt.Errorf("Database save: field %s on scheme %s is required", field.Name, mDat.Scheme.Entity())
			}
			continue
		}

		if field.RefKind == 0 {
			// is not a reference

			operations[i].Field = field.Name
			operations[i].Value, err = field.Type.Value(mField)

			v = runFieldTesters(v, field.Name, operations[i].Value, field.Testers...)

			i++
			// some thing went wrong stop
			if err != nil || v != nil && v.Done().Bad() {
				return v, err
			}

		} else if field.RefKind == scheme.RefParent {

			// gets the model as ref
			fmRef, _isOK := getRefModel(mField)
			if !_isOK {
				err = fmt.Errorf(
					"Database new: field %s from scheme %s can't be saved using type %s ref field don't implement model type",
					field.Name, mDat.Scheme.Entity(), mRef.Type().String(),
				)
				return v, err
			}

			var fm Model
			if !fmRef.IsNil() {
				fm = fmRef.Interface().(Model)
			}

			fmDat := model.GetModelData(fm)

			if mpDat != nil && mpField == field.RefField {
				if mpDat.Scheme != field.RefScheme {
					return v, fmt.Errorf("Database save: todo message here")
				}

				if mField.CanSet() {
					updateParentModel(mpDat.Scheme, reflect.ValueOf(mp), fmRef)
				}

				if fmDat != nil {
					*fmDat = *mpDat // this represents the same model
				} else {
					fmDat = mpDat
				}

			} else if fmDat != nil {
				switch fmDat.Mark {
				case model.WasCreated, model.WasLoaded, model.WasUpdated:
				case model.WasRemoved:
					fmDat.Key = ""
				default:
					return v, fmt.Errorf("Databse save: todo parent field need to be saved before")
				}
			} else {
				// ref model has the field but is nil, skip
				if field.Required {
					return v, fmt.Errorf("Database save: field %s on scheme %s is required", field.Name, mDat.Scheme.Entity())
				}
				continue
			}

			if fmDat.Key == "" {
				if field.Required {
					return v, fmt.Errorf("Database save: field %s on scheme %s is required", field.Name, mDat.Scheme.Entity())
				}
				continue
			}
			operations[i].Field = field.RefField
			operations[i].Value = reflect.ValueOf(&fmDat.Key).Elem()
			i++

		} else if field.RefKind == scheme.RefChild ||
			field.RefKind == scheme.RefChildren ||
			field.RefKind == scheme.RefRelatesTo ||
			field.RefKind == scheme.RefRelatesBack {
			continue
		} else {
			panic(errors.New("Database new: invalid field kind"))
		}
	}

	switch mDat.Mark {
	case 0:
		mDat.Key, err = db.driver.New(mDat.Scheme, mDat.Scheme.KeyField(), operations[:i]...)
		mDat.Mark = model.WasCreated
	case model.WasCreated, model.WasLoaded, model.WasUpdated:
		_, err = db.driver.Update(mDat.Scheme, mDat.Scheme.KeyField(), mDat.Key, operations[:i]...)
		mDat.Mark = model.WasUpdated
	default:
		err = fmt.Errorf("unexpected mark typ %v", mDat.Mark)
	}

	return v, err
}

func updateParentModel(scheme *scheme.Scheme, pmRef reflect.Value, fmRef reflect.Value) {

	if pmRef.Type().AssignableTo(fmRef.Type()) {
		fmRef.Set(pmRef)
		return
	} else if pmRef.Kind() == reflect.Ptr && fmRef.Kind() == reflect.Struct {
		pmRef = pmRef.Elem()
		if pmRef.Type().AssignableTo(fmRef.Type()) {
			fmRef.Set(pmRef)
			return
		}
	} else if pmRef.Kind() == reflect.Struct && fmRef.Kind() == reflect.Ptr {
		if pmRef.Type().AssignableTo(fmRef.Type().Elem()) {
			if pmRef.CanAddr() {
				fmRef.Set(pmRef.Addr())
			} else {
				if !fmRef.IsNil() {
					fmRef.Elem().Set(pmRef)
				}
			}
			return
		}
	}

	if fmRef.Kind() == reflect.Ptr {
		if fmRef.IsNil() {
			return
		}
		fmRef = fmRef.Elem()
	}

	pmRef.FieldByNameFunc(func(field string) bool {
		fmF := fmRef.FieldByName(field)
		pmF := pmRef.FieldByName(field)

		if fmF.IsValid() && fmF.CanSet() {
			fmF.Set(pmF)
		}

		return false
	})
}

func (db *DB) updateModelRefs(m Model, mRef reflect.Value, v *validation.Validator) (*validation.Validator, error) {

	var err error
	var ok bool

	mDat := model.GetModelData(m)

	for _, field := range mDat.Scheme.Fields() {
		if field.RefKind > 0 {

			mField := mRef.FieldByName(field.Name)
			if !mField.IsValid() {
				continue
			}

			if field.RefKind == scheme.RefChild {
				if mField, ok = getRefModel(mField); ok {
					fm := mField.Interface().(Model)

					fmDat := model.GetModelData(fm)
					fmDat.Scheme = field.RefScheme
					println(field.RefScheme.Entity())
					v, err = db.executeSave(fm, m, field.RefField, mField, v)
					if err != nil || v != nil && v.Done().Bad() {
						return v, err
					}
				} else {
					return v, fmt.Errorf(
						"Database new: field %s from scheme %s can't be saved using type %s ref field don't implement model type",
						field.Name, mDat.Scheme.Entity(), mRef.Type().String(),
					)
				}
			} else if field.RefKind == scheme.RefChildren {
				if validModelSlice(mField.Type()) {

					numOfChild := mField.Len()
					for i := 0; i < numOfChild; i++ {
						fmRef := mField.Index(i)

						if fmRef.IsNil() {
							continue
						}

						fm := fmRef.Interface().(Model)

						fmDat := model.GetModelData(fm)
						fmDat.Scheme = field.RefScheme

						v, err = db.executeSave(fm, m, field.RefField, fmRef, v)
						if err != nil || v != nil && v.Done().Bad() {
							return v, err
						}
					}

				} else {
					return v, fmt.Errorf(
						"Database new: field %s from scheme %s can't be saved using type %s ref field don't implement model type",
						field.Name, mDat.Scheme.Entity(), mRef.Type().String(),
					)
				}
			} else if field.RefKind == scheme.RefRelatesTo {
				panic(errors.New("Not implemented yet"))
			}
		}
	}

	return v, err
}

func (db *DB) executeSave(m, mp Model, mpField string, mRef reflect.Value, v *validation.Validator) (*validation.Validator, error) {

	if mRef.Kind() == reflect.Ptr {
		mRef = mRef.Elem()
	}

	v, err := db.updateModelFields(m, mp, mpField, mRef, v)
	if err == nil && v == nil || v.Done().Good() {
		v, err = db.updateModelRefs(m, mRef, v)
	}

	return v, err
}
