package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"sync"
	"time"
	"webview_demo/model"
)

var rawDB *sql.DB
var handle *model.Queries

func InitDB(dsn string) error {
	var err error
	rawDB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	boil.SetDB(rawDB)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	handle, err = model.Prepare(ctx, rawDB)
	if err != nil {
		return err
	}
	return nil
}

type DBStorage struct {
	MetaID int64
	ctx    context.Context
	mu     sync.Mutex
}

func NewDBStorage() *DBStorage {
	return &DBStorage{
		ctx: context.Background(),
	}
}

func saveMeta(ctx context.Context, q *model.Queries, pc *PostContent) (int64, error) {
	id, err := q.CreateImageSetMeta(ctx, model.CreateImageSetMetaParams{
		Title: pc.Title, OriginUrl: pc.Url, Hash: pc.HashHex(),
	})
	return id, err
}

func saveImg(ctx context.Context, q *model.Queries, refId int64, pi *PostImage) error {
	_, err := q.CreateImageItem(ctx, model.CreateImageItemParams{
		RefMeta: refId, Data: pi.Data, Hash: pi.HashHex(), ContentType: pi.ContentType, Url: pi.Url,
	})
	return err
}

func saveFailed(ctx context.Context, q *model.Queries, refId int64, item string) error {
	_, err := q.CreateImageItemFailed(ctx, model.CreateImageItemFailedParams{
		RefMeta: refId, FailedUrl: item, Hash: HashHex([]byte(item)),
	})
	return err
}

func getMetaId(ctx context.Context, q *model.Queries, pc *PostContent) (int64, error) {
	id, err := q.GetImageSetMetaByHash(ctx, pc.HashHex())
	return id, err
}

func getImgId(ctx context.Context, q *model.Queries, pi *PostImage) (int64, error) {
	id, err := q.GetImageItemByHash(ctx, pi.HashHex())
	return id, err
}

func getFailedId(ctx context.Context, q *model.Queries, url string) (int64, error) {
	id, err := q.GetImageFailedByHash(ctx, HashHex([]byte(url)))
	return id, err
}

func (d *DBStorage) createIfNotExist(pc *PostContent) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	ctx, cancel := context.WithTimeout(d.ctx, 5*time.Second)
	defer cancel()

	if d.MetaID == 0 {
		id, err := getMetaId(ctx, handle, pc)
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				if id, err = saveMeta(ctx, handle, pc); err != nil {
					return err
				} else {
					d.MetaID = id
				}

			} else {
				return err
			}

		} else {
			d.MetaID = id
		}
	}
	return nil
}

func (d *DBStorage) Save(pc *PostContent, pi *PostImage) error {
	if err := d.createIfNotExist(pc); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(d.ctx, 5*time.Second)
	defer cancel()

	_, err := getImgId(ctx, handle, pi)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	return saveImg(ctx, handle, d.MetaID, pi)
}

func (d *DBStorage) SaveFailed(pc *PostContent, url string) error {
	if err := d.createIfNotExist(pc); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(d.ctx, 5*time.Second)
	defer cancel()

	_, err := getFailedId(ctx, handle, url)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return saveFailed(ctx, handle, d.MetaID, url)
}

func (d *DBStorage) MkdirAll(pc *PostContent) error {
	return d.createIfNotExist(pc)
}

func (d *DBStorage) tx(ctx context.Context, f func(q *model.Queries) error) error {
	tx, err := rawDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := handle.WithTx(tx)
	err = f(q)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			return fmt.Errorf("%v, rollback %v", err, rErr)
		}
		return err
	}
	return tx.Commit()
}
