package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/middlewares"
	"github.com/Ord1nI/metrix/internal/repo"
	"github.com/Ord1nI/metrix/internal/repo/database"
	"github.com/Ord1nI/metrix/internal/repo/storage"
    _ "net/http/pprof"
)

type Server struct {
	Router      chi.Router
	Middlewares chi.Middlewares
	Config      Config
	Repo        repo.Repo
	Logger      logger.Logger
}

func New() (*Server, error) {
	Logger, err := logger.New()
	if err != nil {
		return nil, err
	}
	s := Server{
		Logger: Logger,
	}
	s.Logger.Infoln("Logger inited successfuly")

	err = s.GetConf()

	if err != nil {
		s.Logger.Errorln("error while gettin conf")
		return nil, err
	}
	s.Logger.Infoln("succesfuly getting conf", s.Config)

    s.Add(middlewares.HeadMW(s.Logger))

	return &s, nil
}

func (s *Server) Init() error {
	err := s.InitRepo()
	if err != nil {
		return err
	}
	s.Logger.Infoln("Repo inited successfuly")

	s.InitRouter(s.Middlewares...)
	return nil
}

func (s *Server) Add(mw ...func(http.Handler) http.Handler) error {
	s.Middlewares = append(s.Middlewares, mw...)
	return nil
}

func(s *Server) RunProff(addres string) {
    go http.ListenAndServe(addres, nil)
}

func (s *Server) Run() error {
	err := s.Init()
	if err != nil {
		s.Logger.Errorln("Fail while starting server")
		return err
	}
	defer s.Repo.Close()
	if s.Router != nil {
		http.ListenAndServe(s.Config.Address, s.Router)
	}
	return errors.New("router not initialized")
}

func (s *Server) InitRepo() error {
	var errM error
	if s.Config.DBdsn != "" {
		err := s.initDB()
		if err == nil {
			return nil
		} else {
			s.Logger.Infoln("Error while connecting to database")
		}
	}

	err := s.initStor()

	if err != nil {
		errM = errors.Join(errM, err)
	}
	return errM
}

func (s *Server) initDB() error {
	s.Logger.Infoln("Trying connection to database")
	db, err := database.NewDB(s.Config.DBdsn, time.Millisecond*500)
	if err != nil {
		s.Logger.Infoln("Failed connecting to database")
		return err
	}

	err = db.Ping()
	if err != nil {
		s.Logger.Infoln("Failed connecting to database")
		return err
	}

	err = db.CreateTable()
	if err != nil {
		s.Logger.Infoln("Failed connecting to database")
		return err
	}
	s.Repo = db
	s.Logger.Infoln("Succesfully connected to database")
	return nil
}

func (s *Server) initStor() error {
	stor := storage.NewMemStorage()

	s.Repo = stor

	if s.Config.FileStoragePath != "" {
		if s.Config.Restore {
			s.Logger.Infoln("Trying to restore data from file")
			err := stor.GetFromFile(s.Config.FileStoragePath)
			if err != nil {
				s.Logger.Infoln("Error while restoring data from file")
			} else {
				s.Logger.Infoln("Data restored from file")
			}
		}

		if s.Config.StoreInterval != 0 {

			s.Logger.Infoln("Starting Data saver with interval: ", s.Config.StoreInterval)
			go stor.StartDataSaver(s.Config.StoreInterval, s.Config.FileStoragePath)
			s.Logger.Infoln("Data saver started successfuly")
		} else {
			s.Logger.Infoln("Starting immediate Data saver")
			s.Add(middlewares.FileWriterWM(s.Logger, stor, s.Config.FileStoragePath))
		}
		return nil
	}

	s.Logger.Infoln("Storage inited with out file saving")
	return nil
}
