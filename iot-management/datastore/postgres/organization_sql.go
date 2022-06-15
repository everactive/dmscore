// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Management Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package postgres

const organizationUserAccessSQL = `
SELECT EXISTS(
	select id from organization_user
	where org_id=$1 and username=$2
)
`

const listUserOrganizationsSQL = `
	select a.code, a.name
	from organization a
	inner join organization_user l on a.code = l.org_id
	where l.username=$1
`
const listOrganizationsSQL = `
	select code, name
	from organization
	where $1 = $1
`

const updateOrganizationSQL = `
	update organization
	set name=$2
	where code = $1
`

const deleteOrganizationUserAccessSQL = `
	delete from organization_user
	where org_id=$1 and username=$2
`

const createOrganizationUserAccessSQL = `
	insert into organization_user (org_id, username)
	values ($1, $2)
`
