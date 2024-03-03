package core

import (
	"github.com/orangeseeds/blitzbase/refactored/store"
	upper "github.com/orangeseeds/blitzbase/store"
)

type App interface {
	Publisher() *upper.Publisher
	Store() store.Store
	Addr() string

	OnStart() *Hook[*AppEvent]
	OnTerminate() *Hook[*AppEvent]

	OnAdminIndex() *Hook[*AdminEvent]
	OnAdminDetail() *Hook[*AdminEvent]
	OnAdminCreate() *Hook[*AdminEvent]
	OnAdminUpdate() *Hook[*AdminEvent]
	OnAdminDelete() *Hook[*AdminEvent]
	OnAdminAuth() *Hook[*AdminEvent]

	OnCollectionIndex() *Hook[*CollectionEvent]
	OnCollectionDetail() *Hook[*CollectionEvent]
	OnCollectionCreate() *Hook[*CollectionEvent]
	OnCollectionUpdate() *Hook[*CollectionEvent]
	OnCollectionDelete() *Hook[*CollectionEvent]

	OnRecordIndex() *Hook[*RecordEvent]
	OnRecordDetail() *Hook[*RecordEvent]
	OnRecordCreate() *Hook[*RecordEvent]
	OnRecordUpdate() *Hook[*RecordEvent]
	OnRecordDelete() *Hook[*RecordEvent]
	OnRecordAuth() *Hook[*RecordEvent]
}

type DBAppConfig struct {
	DbPath     string
	ServerAddr string
}

type DBApp struct {
	publisher   *upper.Publisher
	store       store.Store
	dbDir       string
	addr        string
	onStart     *Hook[*AppEvent]
	onTerminate *Hook[*AppEvent]

	onAdminIndex  *Hook[*AdminEvent]
	onAdminDetail *Hook[*AdminEvent]
	onAdminCreate *Hook[*AdminEvent]
	onAdminUpdate *Hook[*AdminEvent]
	onAdminDelete *Hook[*AdminEvent]
	onAdminAuth   *Hook[*AdminEvent]

	onCollectionIndex  *Hook[*CollectionEvent]
	onCollectionDetail *Hook[*CollectionEvent]
	onCollectionCreate *Hook[*CollectionEvent]
	onCollectionUpdate *Hook[*CollectionEvent]
	onCollectionDelete *Hook[*CollectionEvent]

	onRecordIndex  *Hook[*RecordEvent]
	onRecordDetail *Hook[*RecordEvent]
	onRecordCreate *Hook[*RecordEvent]
	onRecordUpdate *Hook[*RecordEvent]
	onRecordDelete *Hook[*RecordEvent]
	onRecordAuth   *Hook[*RecordEvent]
}

func NewDBApp(config DBAppConfig, store store.Store) *DBApp {
	return &DBApp{
		dbDir:     config.DbPath,
		addr:      config.ServerAddr,
		store:     store,
		publisher: upper.NewPublisher(),

		onStart:     &Hook[*AppEvent]{},
		onTerminate: &Hook[*AppEvent]{},

		onCollectionIndex:  &Hook[*CollectionEvent]{},
		onCollectionDetail: &Hook[*CollectionEvent]{},
		onCollectionCreate: &Hook[*CollectionEvent]{},
		onCollectionUpdate: &Hook[*CollectionEvent]{},
		onCollectionDelete: &Hook[*CollectionEvent]{},

		onAdminIndex:  &Hook[*AdminEvent]{},
		onAdminCreate: &Hook[*AdminEvent]{},
		onAdminAuth:   &Hook[*AdminEvent]{},
		onAdminDetail: &Hook[*AdminEvent]{},
		onAdminUpdate: &Hook[*AdminEvent]{},
		onAdminDelete: &Hook[*AdminEvent]{},

		onRecordIndex:  &Hook[*RecordEvent]{},
		onRecordDetail: &Hook[*RecordEvent]{},
		onRecordCreate: &Hook[*RecordEvent]{},
		onRecordUpdate: &Hook[*RecordEvent]{},
		onRecordDelete: &Hook[*RecordEvent]{},
		onRecordAuth:   &Hook[*RecordEvent]{},
	}
}

func (a *DBApp) Start(addr string) {
	a.addr = addr
	a.onStart.Trigger(&AppEvent{
		Type: StartEvent,
		App:  a,
	})
}
func (a *DBApp) Publisher() *upper.Publisher { return a.publisher }

func (a *DBApp) OnStart() *Hook[*AppEvent]     { return a.onStart }
func (a *DBApp) OnTerminate() *Hook[*AppEvent] { return a.onTerminate }

func (a *DBApp) OnAdminIndex() *Hook[*AdminEvent]  { return a.onAdminIndex }
func (a *DBApp) OnAdminDetail() *Hook[*AdminEvent] { return a.onAdminDetail }
func (a *DBApp) OnAdminCreate() *Hook[*AdminEvent] { return a.onAdminCreate }
func (a *DBApp) OnAdminUpdate() *Hook[*AdminEvent] { return a.onAdminUpdate }
func (a *DBApp) OnAdminDelete() *Hook[*AdminEvent] { return a.onAdminDelete }
func (a *DBApp) OnAdminAuth() *Hook[*AdminEvent]   { return a.onAdminAuth }

func (a *DBApp) OnCollectionIndex() *Hook[*CollectionEvent]  { return a.onCollectionIndex }
func (a *DBApp) OnCollectionDetail() *Hook[*CollectionEvent] { return a.onCollectionDetail }
func (a *DBApp) OnCollectionCreate() *Hook[*CollectionEvent] { return a.onCollectionCreate }
func (a *DBApp) OnCollectionUpdate() *Hook[*CollectionEvent] { return a.onCollectionUpdate }
func (a *DBApp) OnCollectionDelete() *Hook[*CollectionEvent] { return a.onCollectionDelete }

func (a *DBApp) OnRecordIndex() *Hook[*RecordEvent]  { return a.onRecordIndex }
func (a *DBApp) OnRecordDetail() *Hook[*RecordEvent] { return a.onRecordDetail }
func (a *DBApp) OnRecordCreate() *Hook[*RecordEvent] { return a.onRecordCreate }
func (a *DBApp) OnRecordUpdate() *Hook[*RecordEvent] { return a.onRecordUpdate }
func (a *DBApp) OnRecordDelete() *Hook[*RecordEvent] { return a.onRecordDelete }
func (a *DBApp) OnRecordAuth() *Hook[*RecordEvent]   { return a.onRecordAuth }

func (a *DBApp) Store() store.Store { return a.store }

func (a *DBApp) Addr() string { return a.addr }
