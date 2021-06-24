package core

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"

	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

var (
	PaxosBucket = []byte("paxos")
	TermKey     = []byte("term")
	FirstKey    = []byte("first")
	ValueKey    = []byte("value")
	VoteKey     = []byte("vote")

	EmptyTerm  = Term{}
	EmptyFirst = Value("")
	EmptyValue = Value("")
	EmptyVote  = Vote{}
)

type Storage struct {
	Path string
	DB   *bolt.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}

	storage := Storage{
		Path: path,
		DB:   db,
	}
	return &storage, nil
}

func (storage *Storage) Create(term Term, first Value) error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		paxosBucket, err := tx.CreateBucketIfNotExists(PaxosBucket)
		if err != nil {
			return err
		}

		bs := paxosBucket.Get(TermKey)
		if len(bs) != 0 {
			return nil
		}
		bs, err = jsoniter.Marshal(term)
		if err != nil {
			return err
		}
		err = paxosBucket.Put(TermKey, bs)
		if err != nil {
			return err
		}

		if first.IsNull() {
			return nil
		}
		bs = paxosBucket.Get(FirstKey)
		if len(bs) != 0 {
			return nil
		}
		bs, err = jsoniter.Marshal(first)
		if err != nil {
			return err
		}
		return paxosBucket.Put(FirstKey, bs)
	})
}

func (storage *Storage) Put(key []byte, value []byte) error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		paxosBucket := tx.Bucket(PaxosBucket)
		return paxosBucket.Put(key, value)
	})
}

func (storage *Storage) Get(key []byte) ([]byte, error) {
	var value []byte
	err := storage.DB.View(func(tx *bolt.Tx) error {
		paxosBucket := tx.Bucket(PaxosBucket)
		value = paxosBucket.Get(key)
		return nil
	})
	return value, err
}

func (storage *Storage) PutTerm(term Term) error {
	bs, err := jsoniter.Marshal(term)
	if err != nil {
		return err
	}
	return storage.Put(TermKey, bs)
}

func (storage *Storage) GetTerm() (Term, error) {
	bs, err := storage.Get(TermKey)
	if len(bs) == 0 || err != nil {
		return EmptyTerm, err
	}

	var term Term
	err = jsoniter.Unmarshal(bs, &term)
	return term, err
}

func (storage *Storage) PutFirst(first Value) error {
	bs, err := jsoniter.Marshal(first)
	if err != nil {
		return err
	}
	return storage.Put(FirstKey, bs)
}

func (storage *Storage) GetFirst() (Value, error) {
	bs, err := storage.Get(FirstKey)
	if len(bs) == 0 || err != nil {
		return EmptyFirst, err
	}

	var first Value
	err = jsoniter.Unmarshal(bs, &first)
	return first, err
}

func (storage *Storage) PutValue(value Value) error {
	bs, err := jsoniter.Marshal(value)
	if err != nil {
		return err
	}
	return storage.Put(ValueKey, bs)
}

func (storage *Storage) GetValue() (Value, error) {
	bs, err := storage.Get(ValueKey)
	if len(bs) == 0 || err != nil {
		return EmptyValue, err
	}

	var value Value
	err = jsoniter.Unmarshal(bs, &value)
	return value, err
}

func (storage *Storage) PutVote(vote Vote) error {
	bs, err := jsoniter.Marshal(vote)
	if err != nil {
		return err
	}
	return storage.Put(VoteKey, bs)
}

func (storage *Storage) GetVote() (Vote, error) {
	bs, err := storage.Get(VoteKey)
	if len(bs) == 0 || err != nil {
		return EmptyVote, err
	}

	var vote Vote
	err = jsoniter.Unmarshal(bs, &vote)
	return vote, err
}

func (storage *Storage) Delete() error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(PaxosBucket)
	})
}

func (storage *Storage) Stop() {
	err := storage.DB.Close()
	if err != nil {
		log.Warn().Err(err).Msg("storage stop with error")
	}
}
