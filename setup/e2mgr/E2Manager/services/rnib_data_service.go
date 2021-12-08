//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package services

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"net"
	"time"
)

type RNibDataService interface {
	SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error
	UpdateNodebInfo(nodebInfo *entities.NodebInfo) error
	SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error
	GetNodeb(ranName string) (*entities.NodebInfo, error)
	GetListNodebIds() ([]*entities.NbIdentity, error)
	PingRnib() bool
	GetE2TInstance(address string) (*entities.E2TInstance, error)
	GetE2TInstances(addresses []string) ([]*entities.E2TInstance, error)
	GetE2TAddresses() ([]string, error)
	SaveE2TInstance(e2tInstance *entities.E2TInstance) error
	SaveE2TAddresses(addresses []string) error
	GetE2TInstanceNoLogs(address string) (*entities.E2TInstance, error)
	GetE2TInstancesNoLogs(addresses []string) ([]*entities.E2TInstance, error)
	SaveE2TInstanceNoLogs(e2tInstance *entities.E2TInstance) error
	GetE2TAddressesNoLogs() ([]string, error)
	RemoveE2TInstance(e2tAddress string) error
	UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error
	RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error
}

type rNibDataService struct {
	logger        *logger.Logger
	rnibReader    reader.RNibReader
	rnibWriter    rNibWriter.RNibWriter
	maxAttempts   int
	retryInterval time.Duration
}

func NewRnibDataService(logger *logger.Logger, config *configuration.Configuration, rnibReader reader.RNibReader, rnibWriter rNibWriter.RNibWriter) *rNibDataService {
	return &rNibDataService{
		logger:        logger,
		rnibReader:    rnibReader,
		rnibWriter:    rnibWriter,
		maxAttempts:   config.MaxRnibConnectionAttempts,
		retryInterval: time.Duration(config.RnibRetryIntervalMs) * time.Millisecond,
	}
}

func (w *rNibDataService) RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error {
	err := w.retry("RemoveServedNrCells", func() (err error) {
		err = w.rnibWriter.RemoveServedNrCells(inventoryName, servedNrCells)
		return
	})

	return err
}

func (w *rNibDataService) UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error {
	w.logger.Infof("#RnibDataService.UpdateGnbCells - nodebInfo: %s, servedNrCells: %s", nodebInfo, servedNrCells)

	err := w.retry("UpdateGnbCells", func() (err error) {
		err = w.rnibWriter.UpdateGnbCells(nodebInfo, servedNrCells)
		return
	})

	return err
}

func (w *rNibDataService) UpdateNodebInfo(nodebInfo *entities.NodebInfo) error {
	w.logger.Infof("#RnibDataService.UpdateNodebInfo - nodebInfo: %s", nodebInfo)

	err := w.retry("UpdateNodebInfo", func() (err error) {
		err = w.rnibWriter.UpdateNodebInfo(nodebInfo)
		return
	})

	return err
}

func (w *rNibDataService) SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error {
	w.logger.Infof("#RnibDataService.SaveNodeb - nbIdentity: %s, nodebInfo: %s", nbIdentity, nb)

	err := w.retry("SaveNodeb", func() (err error) {
		err = w.rnibWriter.SaveNodeb(nbIdentity, nb)
		return
	})

	return err
}

func (w *rNibDataService) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error {
	w.logger.Infof("#RnibDataService.SaveRanLoadInformation - inventoryName: %s, ranLoadInformation: %s", inventoryName, ranLoadInformation)

	err := w.retry("SaveRanLoadInformation", func() (err error) {
		err = w.rnibWriter.SaveRanLoadInformation(inventoryName, ranLoadInformation)
		return
	})

	return err
}

func (w *rNibDataService) GetNodeb(ranName string) (*entities.NodebInfo, error) {

	var nodeb *entities.NodebInfo = nil

	err := w.retry("GetNodeb", func() (err error) {
		nodeb, err = w.rnibReader.GetNodeb(ranName)
		return
	})

	if err == nil {
		w.logger.Infof("#RnibDataService.GetNodeb - RAN name: %s, connection status: %s, associated E2T: %s", nodeb.RanName, nodeb.ConnectionStatus, nodeb.AssociatedE2TInstanceAddress)
	}

	return nodeb, err
}

func (w *rNibDataService) GetListNodebIds() ([]*entities.NbIdentity, error) {
	var nodeIds []*entities.NbIdentity = nil

	err := w.retry("GetListNodebIds", func() (err error) {
		nodeIds, err = w.rnibReader.GetListNodebIds()
		return
	})

	if err == nil {
		w.logger.Infof("#RnibDataService.GetListNodebIds - RANs count: %d", len(nodeIds))
	}

	return nodeIds, err
}

func (w *rNibDataService) GetE2TInstance(address string) (*entities.E2TInstance, error) {
	var e2tInstance *entities.E2TInstance = nil

	err := w.retry("GetE2TInstance", func() (err error) {
		e2tInstance, err = w.rnibReader.GetE2TInstance(address)
		return
	})

	if err == nil {
		w.logger.Infof("#RnibDataService.GetE2TInstance - E2T instance address: %s, state: %s, associated RANs count: %d, keep Alive ts: %d", e2tInstance.Address, e2tInstance.State, len(e2tInstance.AssociatedRanList), e2tInstance.KeepAliveTimestamp)
	}

	return e2tInstance, err
}

func (w *rNibDataService) GetE2TInstanceNoLogs(address string) (*entities.E2TInstance, error) {
	var e2tInstance *entities.E2TInstance = nil

	err := w.retry("GetE2TInstance", func() (err error) {
		e2tInstance, err = w.rnibReader.GetE2TInstance(address)
		return
	})

	return e2tInstance, err
}

func (w *rNibDataService) GetE2TInstances(addresses []string) ([]*entities.E2TInstance, error) {
	w.logger.Infof("#RnibDataService.GetE2TInstances - addresses: %s", addresses)
	var e2tInstances []*entities.E2TInstance = nil

	err := w.retry("GetE2TInstance", func() (err error) {
		e2tInstances, err = w.rnibReader.GetE2TInstances(addresses)
		return
	})

	return e2tInstances, err
}

func (w *rNibDataService) GetE2TInstancesNoLogs(addresses []string) ([]*entities.E2TInstance, error) {

	var e2tInstances []*entities.E2TInstance = nil

	err := w.retry("GetE2TInstance", func() (err error) {
		e2tInstances, err = w.rnibReader.GetE2TInstances(addresses)
		return
	})

	return e2tInstances, err
}

func (w *rNibDataService) GetE2TAddresses() ([]string, error) {

	var e2tAddresses []string = nil

	err := w.retry("GetE2TAddresses", func() (err error) {
		e2tAddresses, err = w.rnibReader.GetE2TAddresses()
		return
	})

	if err == nil {
		w.logger.Infof("#RnibDataService.GetE2TAddresses - addresses: %s", e2tAddresses)
	}

	return e2tAddresses, err
}

func (w *rNibDataService) GetE2TAddressesNoLogs() ([]string, error) {

	var e2tAddresses []string = nil

	err := w.retry("GetE2TAddresses", func() (err error) {
		e2tAddresses, err = w.rnibReader.GetE2TAddresses()
		return
	})

	return e2tAddresses, err
}

func (w *rNibDataService) SaveE2TInstance(e2tInstance *entities.E2TInstance) error {
	w.logger.Infof("#RnibDataService.SaveE2TInstance - E2T instance address: %s, podName: %s, state: %s, associated RANs count: %d, keep Alive ts: %d", e2tInstance.Address, e2tInstance.PodName, e2tInstance.State, len(e2tInstance.AssociatedRanList), e2tInstance.KeepAliveTimestamp)

	return w.SaveE2TInstanceNoLogs(e2tInstance)
}

func (w *rNibDataService) SaveE2TInstanceNoLogs(e2tInstance *entities.E2TInstance) error {

	err := w.retry("SaveE2TInstance", func() (err error) {
		err = w.rnibWriter.SaveE2TInstance(e2tInstance)
		return
	})

	return err
}

func (w *rNibDataService) SaveE2TAddresses(addresses []string) error {
	w.logger.Infof("#RnibDataService.SaveE2TAddresses - addresses: %s", addresses)

	err := w.retry("SaveE2TAddresses", func() (err error) {
		err = w.rnibWriter.SaveE2TAddresses(addresses)
		return
	})

	return err
}

func (w *rNibDataService) RemoveE2TInstance(e2tAddress string) error {
	w.logger.Infof("#RnibDataService.RemoveE2TInstance - e2tAddress: %s", e2tAddress)

	err := w.retry("RemoveE2TInstance", func() (err error) {
		err = w.rnibWriter.RemoveE2TInstance(e2tAddress)
		return
	})

	return err
}

func (w *rNibDataService) PingRnib() bool {
	err := w.retry("GetListNodebIds", func() (err error) {
		_, err = w.rnibReader.GetListNodebIds()
		return
	})

	return !isRnibConnectionError(err)
}

func (w *rNibDataService) retry(rnibFunc string, f func() error) (err error) {
	attempts := w.maxAttempts

	for i := 1; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if !isRnibConnectionError(err) {
			return err
		}
		if i >= attempts {
			w.logger.Errorf("#RnibDataService.retry - after %d attempts of %s, last error: %s", attempts, rnibFunc, err)
			return err
		}
		time.Sleep(w.retryInterval)

		w.logger.Infof("#RnibDataService.retry - retrying %d %s after error: %s", i, rnibFunc, err)
	}
}

func isRnibConnectionError(err error) bool {
	internalErr, ok := err.(*common.InternalError)
	if !ok {
		return false
	}
	_, ok = internalErr.Err.(*net.OpError)

	return ok
}
