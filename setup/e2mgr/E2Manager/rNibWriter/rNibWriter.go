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

package rNibWriter

import (
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/proto"
)

const E2TAddressesKey = "E2TAddresses"

type rNibWriterInstance struct {
	sdl common.ISdlInstance
}

/*
RNibWriter interface allows saving data to the redis DB
*/
type RNibWriter interface {
	SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error
	UpdateNodebInfo(nodebInfo *entities.NodebInfo) error
	SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error
	SaveE2TInstance(e2tInstance *entities.E2TInstance) error
	SaveE2TAddresses(addresses []string) error
	RemoveE2TInstance(e2tAddress string) error
	UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error
	RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error
}

/*
GetRNibWriter returns reference to RNibWriter
*/

func GetRNibWriter(sdl common.ISdlInstance) RNibWriter {
	return &rNibWriterInstance{sdl: sdl}
}


func (w *rNibWriterInstance) RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error {
	cellKeysToRemove := buildCellKeysToRemove(inventoryName, servedNrCells)
	err := w.sdl.Remove(cellKeysToRemove)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

/*
SaveNodeb saves nodeB entity data in the redis DB according to the specified data model
*/
func (w *rNibWriterInstance) SaveNodeb(nbIdentity *entities.NbIdentity, entity *entities.NodebInfo) error {
	isNotEmptyIdentity := isNotEmpty(nbIdentity)

	if isNotEmptyIdentity && entity.GetNodeType() == entities.Node_UNKNOWN {
		return common.NewValidationError(fmt.Sprintf("#rNibWriter.saveNodeB - Unknown responding node type, entity: %v", entity))
	}
	data, err := proto.Marshal(entity)
	if err != nil {
		return common.NewInternalError(err)
	}
	var pairs []interface{}
	key, rNibErr := common.ValidateAndBuildNodeBNameKey(nbIdentity.InventoryName)
	if rNibErr != nil {
		return rNibErr
	}
	pairs = append(pairs, key, data)

	if isNotEmptyIdentity {
		key, rNibErr = common.ValidateAndBuildNodeBIdKey(entity.GetNodeType().String(), nbIdentity.GlobalNbId.GetPlmnId(), nbIdentity.GlobalNbId.GetNbId())
		if rNibErr != nil {
			return rNibErr
		}
		pairs = append(pairs, key, data)
	}

	if entity.GetEnb() != nil {
		pairs, rNibErr = appendEnbCells(nbIdentity.InventoryName, entity.GetEnb().GetServedCells(), pairs)
		if rNibErr != nil {
			return rNibErr
		}
	}
	if entity.GetGnb() != nil {
		pairs, rNibErr = appendGnbCells(nbIdentity.InventoryName, entity.GetGnb().GetServedNrCells(), pairs)
		if rNibErr != nil {
			return rNibErr
		}
	}
	err = w.sdl.Set(pairs)
	if err != nil {
		return common.NewInternalError(err)
	}

	ranNameIdentity := &entities.NbIdentity{InventoryName: nbIdentity.InventoryName}

	if isNotEmptyIdentity {
		nbIdData, err := proto.Marshal(ranNameIdentity)
		if err != nil {
			return common.NewInternalError(err)
		}
		err = w.sdl.RemoveMember(entities.Node_UNKNOWN.String(), nbIdData)
		if err != nil {
			return common.NewInternalError(err)
		}
	} else {
		nbIdentity = ranNameIdentity
	}

	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		return common.NewInternalError(err)
	}
	err = w.sdl.AddMember(entity.GetNodeType().String(), nbIdData)
	if err != nil {
		return common.NewInternalError(err)
	}
	return nil
}

func (w *rNibWriterInstance) UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error {

	pairs, err := buildUpdateNodebInfoPairs(nodebInfo)

	if err != nil {
		return err
	}

	pairs, err = appendGnbCells(nodebInfo.RanName, servedNrCells, pairs)

	if err != nil {
		return err
	}

	err = w.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func buildCellKeysToRemove(inventoryName string, servedNrCellsToRemove []*entities.ServedNRCell) []string {

	cellKeysToRemove := []string{}

	for _, cell := range servedNrCellsToRemove {

		key, _ := common.ValidateAndBuildNrCellIdKey(cell.GetServedNrCellInformation().GetCellId())

		if len(key) != 0 {
			cellKeysToRemove = append(cellKeysToRemove, key)
		}

		key, _ = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetServedNrCellInformation().GetNrPci())

		if len(key) != 0 {
			cellKeysToRemove = append(cellKeysToRemove, key)
		}
	}

	return cellKeysToRemove
}

func buildUpdateNodebInfoPairs(nodebInfo *entities.NodebInfo) ([]interface{}, error) {
	nodebNameKey, rNibErr := common.ValidateAndBuildNodeBNameKey(nodebInfo.GetRanName())

	if rNibErr != nil {
		return []interface{}{}, rNibErr
	}

	nodebIdKey, buildNodebIdKeyError := common.ValidateAndBuildNodeBIdKey(nodebInfo.GetNodeType().String(), nodebInfo.GlobalNbId.GetPlmnId(), nodebInfo.GlobalNbId.GetNbId())

	data, err := proto.Marshal(nodebInfo)

	if err != nil {
		return []interface{}{}, common.NewInternalError(err)
	}

	pairs := []interface{}{nodebNameKey, data}

	if buildNodebIdKeyError == nil {
		pairs = append(pairs, nodebIdKey, data)
	}

	return pairs, nil
}

/*
UpdateNodebInfo...
*/
func (w *rNibWriterInstance) UpdateNodebInfo(nodebInfo *entities.NodebInfo) error {

	pairs, err := buildUpdateNodebInfoPairs(nodebInfo)

	if err != nil {
		return err
	}

	err = w.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

/*
SaveRanLoadInformation stores ran load information for the provided ran
*/
func (w *rNibWriterInstance) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error {

	key, rnibErr := common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if rnibErr != nil {
		return rnibErr
	}

	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = w.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) SaveE2TInstance(e2tInstance *entities.E2TInstance) error {

	key, rnibErr := common.ValidateAndBuildE2TInstanceKey(e2tInstance.Address)

	if rnibErr != nil {
		return rnibErr
	}

	data, err := json.Marshal(e2tInstance)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = w.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) SaveE2TAddresses(addresses []string) error {

	data, err := json.Marshal(addresses)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, E2TAddressesKey, data)

	err = w.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) RemoveE2TInstance(address string) error {
	key, rNibErr := common.ValidateAndBuildE2TInstanceKey(address)
	if rNibErr != nil {
		return rNibErr
	}
	err := w.sdl.Remove([]string{key})

	if err != nil {
		return common.NewInternalError(err)
	}
	return nil
}

/*
Close the writer
*/
func Close() {
	//Nothing to do
}

func appendEnbCells(inventoryName string, cells []*entities.ServedCellInfo, pairs []interface{}) ([]interface{}, error) {
	for _, cell := range cells {
		cellEntity := entities.Cell{Type: entities.Cell_LTE_CELL, Cell: &entities.Cell_ServedCellInfo{ServedCellInfo: cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, common.NewInternalError(err)
		}
		key, rNibErr := common.ValidateAndBuildCellIdKey(cell.GetCellId())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetPci())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func appendGnbCells(inventoryName string, cells []*entities.ServedNRCell, pairs []interface{}) ([]interface{}, error) {
	for _, cell := range cells {
		cellEntity := entities.Cell{Type: entities.Cell_NR_CELL, Cell: &entities.Cell_ServedNrCell{ServedNrCell: cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, common.NewInternalError(err)
		}
		key, rNibErr := common.ValidateAndBuildNrCellIdKey(cell.GetServedNrCellInformation().GetCellId())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetServedNrCellInformation().GetNrPci())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func isNotEmpty(nbIdentity *entities.NbIdentity) bool {
	return nbIdentity.GlobalNbId != nil && nbIdentity.GlobalNbId.PlmnId != "" && nbIdentity.GlobalNbId.NbId != ""
}
