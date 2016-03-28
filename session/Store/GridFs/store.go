// mgoStore uses the mgo mongo driver to store sessions on mongo gridfs
package gridFs

import (
	"github.com/CloudyKit/framework/session"
	"github.com/jhsx/qm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"time"
)

// New returns a new store
func New(db, prefix string, mgoSess func() *mgo.Session) session.Store {
	return &Store{
		db:      db,
		prefix:  prefix,
		session: mgoSess,
	}
}

type sessionCloser struct {
	session *mgo.Session
	*mgo.GridFile
}

func (ss *sessionCloser) Close() (err error) {
	if ss.GridFile != nil {
		if err = ss.GridFile.Close(); err != nil {
			ss.session.Close()
			return
		}
	}
	ss.session.Close()
	return
}

type Store struct {
	session    func() *mgo.Session
	db, prefix string
}

func (sessionStore *Store) gridFs(name string, create bool) (*sessionCloser, error) {
	session := sessionStore.session()
	session.SetMode(mgo.Strong, false)
	gridFs := session.DB(sessionStore.db).GridFS(sessionStore.prefix)
	if create {
		gridFs.Remove(name)
		gridFile, err := gridFs.Create(name)
		return &sessionCloser{session: session, GridFile: gridFile}, err
	}
	gridFile, err := gridFs.Open(name)
	return &sessionCloser{session: session, GridFile: gridFile}, err
}

func (sessionStore *Store) Writer(name string) (writer io.WriteCloser, err error) {
	writer, err = sessionStore.gridFs(name, true)
	return
}

func (sessionStore *Store) Reader(name string) (reader io.ReadCloser, err error) {
	reader, err = sessionStore.gridFs(name, false)
	return
}

func (sessionStore *Store) Remove(name string) (err error) {
	sess := sessionStore.session()
	defer sess.Close()
	return sess.DB(sessionStore.db).GridFS(sessionStore.prefix).Remove(name)
}

func (sessionStore *Store) Gc(before time.Time) {
	sess := sessionStore.session()
	defer sess.Close()

	gridFs := sess.DB(sessionStore.db).GridFS(sessionStore.prefix)

	var fileId struct {
		Id bson.ObjectId `bson:"_id"`
	}

	inter := gridFs.Find(qm.Lt("uploadDate", before)).Iter()
	defer inter.Close()

	for inter.Next(&fileId) {
		gridFs.RemoveId(fileId.Id)
	}

	return
}
