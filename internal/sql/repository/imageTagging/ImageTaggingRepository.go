/*
 * Copyright (c) 2020 Devtron Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package repository

import (
	"github.com/go-pg/pg"
	"time"
)

type ImageTaggingAction int

//this action is only allowed for imageComments
const ActionEdit ImageTaggingAction = 1

const ActionSave ImageTaggingAction = 0
const ActionHardDelete ImageTaggingAction = 3
const ActionSoftDelete ImageTaggingAction = 2

type AuditType int

const TagType AuditType = 0
const CommentType AuditType = 1

type ImageTag struct {
	TableName  struct{} `sql:"release_tags" json:",omitempty"  pg:",discard_unknown_columns"`
	Id         int      `sql:"id" json:"id"`
	TagName    string   `sql:"tag_name" json:"tagName"`
	AppId      int      `sql:"app_id" json:"appId"`
	ArtifactId int      `sql:"artifact_id" json:"artifactId"`
	Active     bool     `sql:"active" json:"active"` //this flag is to check soft delete
}

type ImageComment struct {
	TableName  struct{} `sql:"image_comments" json:",omitempty"  pg:",discard_unknown_columns"`
	Id         int      `sql:"id" json:"id"`
	Comment    string   `sql:"comment" json:"comment"`
	ArtifactId int      `sql:"artifact_id" json:"artifactId"`
	UserId     int      `sql:"user_id" json:"-"` //currently not sending userId in json response
}

type ImageTaggingAudit struct {
	TableName  struct{}           `sql:"release_tags" json:",omitempty"  pg:",discard_unknown_columns"`
	Id         int                `sql:"id"`
	Data       string             `sql:"data"`
	DataType   AuditType          `sql:"data_type"`
	ArtifactId int                `sql:"artifact_id"`
	UpdatedOn  time.Time          `sql:"updated_on"`
	UpdatedBy  int                `sql:"updated_by"`
	Action     ImageTaggingAction `sql:"action"`
}

type ImageTaggingRepository interface {
	SaveAuditLogsInBulk(tx *pg.Tx, imageTaggingAudit []*ImageTaggingAudit) error
	SaveReleaseTagsInBulk(tx *pg.Tx, imageTags []*ImageTag) error
	SaveImageComment(tx *pg.Tx, imageComment *ImageComment) error
	GetTagsByAppId(appId int) ([]ImageTag, error)
	GetTagsByArtifactId(artifactId int) ([]ImageTag, error)
	GetImageComment(artifactId int) (ImageComment, error)
	GetImageCommentsByAppId(appId int) ([]ImageComment, error)
	UpdateReleaseTagInBulk(tx *pg.Tx, imageTags []*ImageTag) error
	UpdateImageComment(tx *pg.Tx, imageComment *ImageComment) error
	DeleteReleaseTagInBulk(tx *pg.Tx, imageTags []*ImageTag) error
	StartTx() (*pg.Tx, error)
	RollbackTx(tx *pg.Tx) error
	CommitTx(tx *pg.Tx) error
}

type ImageTaggingRepositoryImpl struct {
	dbConnection *pg.DB
}

func NewImageTaggingRepositoryImpl(db *pg.DB) *ImageTaggingRepositoryImpl {
	return &ImageTaggingRepositoryImpl{
		dbConnection: db,
	}
}

func (impl *ImageTaggingRepositoryImpl) RollbackTx(tx *pg.Tx) error {
	return tx.Rollback()
}
func (impl *ImageTaggingRepositoryImpl) CommitTx(tx *pg.Tx) error {
	return tx.Commit()
}
func (impl *ImageTaggingRepositoryImpl) StartTx() (*pg.Tx, error) {
	return impl.dbConnection.Begin()
}
func (impl *ImageTaggingRepositoryImpl) SaveAuditLogsInBulk(tx *pg.Tx, imageTaggingAudit []*ImageTaggingAudit) error {
	err := tx.Insert(tx, imageTaggingAudit)
	return err
}
func (impl *ImageTaggingRepositoryImpl) SaveReleaseTagsInBulk(tx *pg.Tx, imageTags []*ImageTag) error {
	err := tx.Insert(imageTags)
	return err
}

func (impl *ImageTaggingRepositoryImpl) SaveImageComment(tx *pg.Tx, imageComment *ImageComment) error {
	err := tx.Insert(imageComment)
	return err
}

func (impl *ImageTaggingRepositoryImpl) GetTagsByAppId(appId int) ([]ImageTag, error) {
	res := make([]ImageTag, 0)
	err := impl.dbConnection.Model(&res).
		Where("app_id=?", appId).
		Select()
	return res, err
}

func (impl *ImageTaggingRepositoryImpl) GetTagsByArtifactId(artifactId int) ([]ImageTag, error) {
	res := make([]ImageTag, 0)
	err := impl.dbConnection.Model(&res).
		Where("artifact_id=?", artifactId).
		Select()
	return res, err
}

func (impl *ImageTaggingRepositoryImpl) GetImageComment(artifactId int) (ImageComment, error) {
	res := ImageComment{}
	err := impl.dbConnection.Model(&res).
		Where("artifact_id=?", artifactId).
		Select()
	return res, err
}

func (impl *ImageTaggingRepositoryImpl) GetImageCommentsByAppId(appId int) ([]ImageComment, error) {
	res := make([]ImageComment, 0)
	err := impl.dbConnection.Model(&res).
		Where("app_id=?", appId).
		Select()
	return res, err
}

//this will update the provided release tag
func (impl *ImageTaggingRepositoryImpl) UpdateReleaseTagInBulk(tx *pg.Tx, imageTags []*ImageTag) error {
	//currently tags are not editable, can only be soft deleted or hard delete
	err := tx.Update(imageTags)
	return err
}
func (impl *ImageTaggingRepositoryImpl) UpdateImageComment(tx *pg.Tx, imageComment *ImageComment) error {
	err := tx.Update(imageComment)
	return err
}

func (impl *ImageTaggingRepositoryImpl) DeleteReleaseTagInBulk(tx *pg.Tx, imageTags []*ImageTag) error {
	err := tx.Delete(imageTags)
	return err
}