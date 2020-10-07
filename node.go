package gosql

import (
	"database/sql"
)

func NewDB(driver string, url string, c ConnConfiger) *sql.DB {

	db, err := sql.Open(driver, url)
	if nil != err {
		panic(err)
	}

	db.SetConnMaxLifetime(c.GetMaxLifetime())
	db.SetMaxIdleConns(c.GetMaxIdleConns())
	db.SetMaxOpenConns(c.GetMaxOpenConns())

	return db
}

func NewMasterNode(driver string, c MasterNodeConfiger) *MasterNode {

	return &MasterNode{
		db: NewDB(driver, c.GetMasterUrl(), c),
	}
}

type MasterNode struct {
	db *sql.DB
}

func (m *MasterNode) GetDB() *sql.DB {
	return m.db
}

func NewReplicasNode(driver string, c ReplicaNodeConfiger, opts ...ReplicaNodeOption) *ReplicaNode {
	var dbs []*sql.DB

	for _, url := range c.GetReplicaUrls() {
		dbs = append(dbs, NewDB(driver, url, c))
	}

	node := &ReplicaNode{
		db: dbs,
	}

	return node.WithOptions(opts...)
}

type ReplicaNodeOption interface {
	apply(*ReplicaNode)
}

type replicaOptionFunc func(*ReplicaNode)

func (f replicaOptionFunc) apply(c *ReplicaNode) {
	f(c)
}

func AddReplication(lb Replication) ReplicaNodeOption {
	return replicaOptionFunc(func(r *ReplicaNode) {
		r.SetReplication(lb)
	})
}

type ReplicaNode struct {
	db []*sql.DB
	lb Replication
}

func (r *ReplicaNode) clone() *ReplicaNode {
	copied := *r
	return &copied
}

func (r *ReplicaNode) WithOptions(opts ...ReplicaNodeOption) *ReplicaNode {
	node := r.clone()
	for _, opt := range opts {
		opt.apply(node)
	}
	return node
}

func (r *ReplicaNode) SetReplication(lb Replication) {
	r.lb = lb
}

func (r *ReplicaNode) GetDB() *sql.DB {
	if total := len(r.db); total == 0 {
		return nil
	} else if total == 1 {
		return r.db[0]
	} else {
		return r.db[r.lb.Replicate(total)]
	}
}
