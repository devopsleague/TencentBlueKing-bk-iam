/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-权限中心(BlueKing-IAM) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package dao

//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"iam/pkg/database"
)

// SubjectTemplateGroup  用户/部门-人员模版-用户组关系表
type SubjectTemplateGroup struct {
	PK         int64     `db:"pk"`
	SubjectPK  int64     `db:"subject_pk"`
	TemplateID int64     `db:"template_id"`
	GroupPK    int64     `db:"group_pk"`
	ExpiredAt  int64     `db:"expired_at"`
	CreatedAt  time.Time `db:"created_at"`
}

type SubjectTemplateGroupManager interface {
	BulkCreateWithTx(tx *sqlx.Tx, relations []SubjectTemplateGroup) error
	BulkUpdateExpiredAtWithTx(tx *sqlx.Tx, relations []SubjectRelation) error
	BulkDeleteWithTx(tx *sqlx.Tx, relations []SubjectTemplateGroup) error
	HasRelationExceptTemplate(subjectPK, groupPK, templateID int64) (bool, error)
	GetExpiredAtBySubjectGroup(subjectPK, groupPK int64) (int64, error)
}

type subjectTemplateGroupManager struct {
	DB *sqlx.DB
}

// NewSubjectTemplateGroupManager New SubjectTemplateGroupManager
func NewSubjectTemplateGroupManager() SubjectTemplateGroupManager {
	return &subjectTemplateGroupManager{
		DB: database.GetDefaultDBClient().DB,
	}
}

// BulkCreateWithTx ...
func (m *subjectTemplateGroupManager) BulkCreateWithTx(tx *sqlx.Tx, relations []SubjectTemplateGroup) error {
	if len(relations) == 0 {
		return nil
	}

	sql := `INSERT INTO subject_template_group (
		subject_pk,
		template_id,
		group_pk,
		expired_at
	) VALUES (:subject_pk,
		:template_id,
		:group_pk,
		:expired_at)`
	return database.SqlxBulkInsertWithTx(tx, sql, relations)
}

// BulkUpdateExpiredAtWithTx ...
func (m *subjectTemplateGroupManager) BulkUpdateExpiredAtWithTx(
	tx *sqlx.Tx,
	relations []SubjectRelation,
) error {
	sql := `UPDATE subject_template_group
		 SET expired_at = :expired_at
		 WHERE subject_pk = :subject_pk AND group_pk = :group_pk`
	return database.SqlxBulkUpdateWithTx(tx, sql, relations)
}

// BulkDeleteWithTx ...
func (m *subjectTemplateGroupManager) BulkDeleteWithTx(tx *sqlx.Tx, relations []SubjectTemplateGroup) error {
	if len(relations) == 0 {
		return nil
	}

	sql := `DELETE FROM subject_template_group WHERE subject_pk = ? AND group_pk = ? AND template_id = ?`
	return database.SqlxBulkUpdateWithTx(tx, sql, relations)
}

func (m *subjectTemplateGroupManager) HasRelationExceptTemplate(subjectPK, groupPK, templateID int64) (bool, error) {
	var pk int64
	query := `SELECT
		pk
		FROM subject_template_group
		WHERE subject_pk = ?
		AND group_pk = ?
		AND template_id != ?
		LIMIT 1`
	err := database.SqlxGet(m.DB, &pk, query, subjectPK, groupPK, templateID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *subjectTemplateGroupManager) GetExpiredAtBySubjectGroup(subjectPK, groupPK int64) (int64, error) {
	var expiredAt int64
	query := `SELECT
		 expired_at
		 FROM subject_template_group
		 WHERE subject_pk = ?
		 AND group_pk = ?
		 LIMIT 1`
	err := database.SqlxGet(m.DB, &expiredAt, query, subjectPK, groupPK)
	return expiredAt, err
}