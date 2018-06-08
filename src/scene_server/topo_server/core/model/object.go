/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	frtypes "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"

	"configcenter/src/scene_server/topo_server/core/types"
	metadata "configcenter/src/source_controller/api/metadata"
)

var _ Object = (*object)(nil)

type object struct {
	obj       meta.Object
	isNew     bool
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *object) IsExists() (bool, error) {

	cond := common.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(cli.params.Header.OwnerID).Field(common.BKObjIDField).Eq(cli.ObjectID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjects(context.Background(), util.Headers{
		Language: cli.params.Header.Language,
		OwnerID:  cli.params.Header.OwnerID,
	}, cond.ToMapStr().ToJSON())

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return false, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.ObjectID, rsp.Message)
		return false, cli.params.Err.Error(rsp.Code)
	}

	// TODO: check the result

	return true, nil
}

func (cli *object) Create() error {

	obj := &metadata.ObjectDes{}

	obj.Creator = cli.obj.Creator
	obj.Description = cli.obj.Description
	obj.IsPaused = cli.obj.IsPaused
	obj.IsPre = cli.obj.IsPre
	obj.ObjCls = cli.obj.ObjCls
	obj.Modifier = cli.obj.Modifier
	obj.ObjectID = cli.obj.ObjectID
	obj.ObjectName = cli.obj.ObjectName
	obj.ObjIcon = cli.obj.ObjIcon
	obj.OwnerID = cli.params.Header.OwnerID

	rsp, err := cli.clientSet.ObjectController().Meta().CreateObject(context.Background(), cli.params.Header, obj)

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.ObjectID, rsp.Message)
		return cli.params.Err.Error(rsp.Code)
	}

	// TODO: check the result

	return nil
}

func (cli *object) Update() error {

	data := common.SetValueToMapStrByTags(cli)

	rsp, err := cli.clientSet.ObjectController().Meta().UpdateObject(context.Background(), cli.ObjectID, util.Headers{
		Language: cli.params.Header.Language,
		OwnerID:  cli.params.Header.OwnerID,
	}, data)

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.ObjectID, rsp.Message)
		return cli.params.Err.Error(rsp.Code)
	}

	return nil
}

func (cli *object) Delete() error {
	return nil
}

func (cli *object) Parse(data frtypes.MapStr) error {

	err := common.SetValueToStructByTags(cli, data)
	if nil != err {
		return err
	}

	if 0 == len(cli.ObjectID) {
		return cli.params.Err.Errorf(common.CCErrCommParamsNeedSet, "bk_obj_id")
	}

	if 0 == len(cli.ObjCls) {
		return cli.params.Err.Errorf(common.CCErrCommParamsNeedSet, "bk_classification_id")
	}

	return err
}

func (cli *object) ToMapStr() (frtypes.MapStr, error) {
	return nil, nil
}

func (cli *object) Save() error {
	dataMapStr := common.SetValueToMapStrByTags(cli)

	if cli.isNew {
		cli.create()
	} else {
		cli.update()
	}
	_ = dataMapStr
	return nil
}

func (cli *object) CreateGroup() Group {
	return &group{
		OwnerID:  cli.OwnerID,
		ObjectID: cli.ObjectID,
	}
}

func (cli *object) CreateAttribute() Attribute {
	return &attribute{
		OwnerID:  cli.OwnerID,
		ObjectID: cli.ObjectID,
	}
}

func (cli *object) GetAttributes() ([]Attribute, error) {
	return nil, nil
}

func (cli *object) GetGroups() ([]Group, error) {
	return nil, nil
}

func (cli *object) SetClassification(class Classification) {
	cli.ObjCls = class.GetID()
}

func (cli *object) GetClassification() (Classification, error) {
	return nil, nil
}

func (cli *object) SetIcon(objectIcon string) {
	cli.ObjIcon = objectIcon
}

func (cli *object) GetIcon() string {
	return cli.ObjIcon
}

func (cli *object) SetID(objectID string) {
	cli.ObjectID = objectID
}

func (cli *object) GetID() string {
	return cli.ObjectID
}

func (cli *object) SetName(objectName string) {
	cli.ObjectName = objectName
}

func (cli *object) GetName() string {
	return cli.ObjectName
}

func (cli *object) SetIsPre(isPre bool) {
	cli.IsPre = isPre
}

func (cli *object) GetIsPre() bool {
	return cli.IsPre
}

func (cli *object) SetIsPaused(isPaused bool) {
	cli.IsPaused = isPaused
}

func (cli *object) GetIsPaused() bool {
	return cli.IsPaused
}

func (cli *object) SetPosition(position string) {
	cli.Position = position
}

func (cli *object) GetPosition() string {
	return cli.Position
}

func (cli *object) SetSupplierAccount(supplierAccount string) {
	cli.OwnerID = supplierAccount
}

func (cli *object) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *object) SetDescription(description string) {
	cli.Description = description
}

func (cli *object) GetDescription() string {
	return cli.Description
}

func (cli *object) SetCreator(creator string) {
	cli.Creator = creator
}

func (cli *object) GetCreator() string {
	return cli.Creator
}

func (cli *object) SetModifier(modifier string) {
	cli.Modifier = modifier
}

func (cli *object) GetModifier() string {
	return cli.Modifier
}
