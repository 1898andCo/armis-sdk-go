// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

type BaseResponse struct {
	Success bool `json:"success"`
}

type CreateRoleAPIResponse struct {
	BaseResponse
}

type DeleteRoleAPIResponse struct {
	BaseResponse
}

type RolesAPIResponse struct {
	Data []RoleSettings `json:"data,omitempty"`
	BaseResponse
}

type RoleSettings struct {
	Name        string      `json:"name"`
	Permissions Permissions `json:"permissions"`
	ID          int         `json:"roleId,omitempty"`
	ViprRole    bool        `json:"viprRole,omitempty"`
}

type Permissions struct {
	AdvancedPermissions AdvancedPermissions `json:"advancedPermissions"`
	Alert               Alert               `json:"alert"`
	Device              Device              `json:"device"`
	Policy              Policy              `json:"policy"`
	Report              ReportPermissions   `json:"report"`
	RiskFactor          RiskFactor          `json:"risk_factor"`
	Settings            Settings            `json:"settings"`
	User                User                `json:"user"`
	Vulnerability       Vulnerability       `json:"vulnerability"`
}

type AdvancedPermissions struct {
	All        bool           `json:"all,omitempty"`
	Behavioral Behavioral     `json:"behavioral"`
	Device     DeviceAdvanced `json:"device"`
}

type Behavioral struct {
	All             bool       `json:"all,omitempty"`
	ApplicationName Permission `json:"applicationName"`
	HostName        Permission `json:"hostName"`
	ServiceName     Permission `json:"serviceName"`
}

type DeviceAdvanced struct {
	All          bool       `json:"all,omitempty"`
	DeviceNames  Permission `json:"deviceNames"`
	IPAddresses  Permission `json:"ipAddresses"`
	MACAddresses Permission `json:"macAddresses"`
	PhoneNumbers Permission `json:"phoneNumbers"`
}

type Permission struct {
	All bool `json:"all,omitempty"`
}

type Alert struct {
	All    bool       `json:"all,omitempty"`
	Manage Manage     `json:"manage"`
	Read   Permission `json:"read"`
}

type Manage struct {
	All              bool       `json:"all,omitempty"`
	Resolve          Permission `json:"resolve"`
	WhitelistDevices Permission `json:"whitelistDevices"`
}

type Device struct {
	All    bool         `json:"all,omitempty"`
	Manage ManageDevice `json:"manage"`
	Read   Permission   `json:"read"`
}

type ManageDevice struct {
	All                bool       `json:"all,omitempty"`
	Create             Permission `json:"create"`
	Delete             Permission `json:"delete"`
	Edit               Permission `json:"edit"`
	Enforce            Enforce    `json:"enforce"`
	Merge              Permission `json:"merge"`
	RequestDeletedData Permission `json:"request_deleted_data"`
	Tags               Permission `json:"tags"`
}

type Enforce struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create"`
	Delete Permission `json:"delete"`
}

type Policy struct {
	All    bool       `json:"all,omitempty"`
	Manage Permission `json:"manage"`
	Read   Permission `json:"read"`
}

type ReportPermissions struct {
	All    bool         `json:"all,omitempty"`
	Export Permission   `json:"export"`
	Manage ManageReport `json:"manage"`
	Read   Permission   `json:"read"`
}

type ManageReport struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create"`
	Delete Permission `json:"delete"`
	Edit   Permission `json:"edit"`
}

type RiskFactor struct {
	All    bool       `json:"all,omitempty"`
	Manage ManageRisk `json:"manage"`
	Read   Permission `json:"read"`
}

type ManageRisk struct {
	All           bool          `json:"all,omitempty"`
	Customization Customization `json:"customization"`
	Status        Status        `json:"status"`
}

type Customization struct {
	All     bool       `json:"all,omitempty"`
	Create  Permission `json:"create"`
	Disable Permission `json:"disable"`
	Edit    Permission `json:"edit"`
}

type Status struct {
	All     bool       `json:"all,omitempty"`
	Ignore  Permission `json:"ignore"`
	Resolve Permission `json:"resolve"`
}

type Settings struct {
	All              bool            `json:"all,omitempty"`
	AuditLog         Permission      `json:"auditLog"`
	Boundary         Boundary        `json:"boundary"`
	BusinessImpact   ManageAndRead   `json:"businessImpact"`
	Collector        ManageAndRead   `json:"collector"`
	CustomProperties ManageAndRead   `json:"customProperties"`
	Integration      ManageAndRead   `json:"integration"`
	InternalIps      ManageAndRead   `json:"internalIps"`
	Notifications    ManageAndRead   `json:"notifications"`
	OIDC             ManageAndRead   `json:"oidc"`
	SAML             ManageAndRead   `json:"saml"`
	SecretKey        Permission      `json:"secretKey"`
	SecuritySettings Permission      `json:"securitySettings"`
	SitesAndSensors  SitesAndSensors `json:"sitesAndSensors"`
	UsersAndRoles    UsersAndRoles   `json:"usersAndRoles"`
}

type Boundary struct {
	All    bool           `json:"all,omitempty"`
	Manage ManageBoundary `json:"manage"`
	Read   Permission     `json:"read"`
}

type ManageBoundary struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create"`
	Delete Permission `json:"delete"`
	Edit   Permission `json:"edit"`
}

type ManageAndRead struct {
	All    bool       `json:"all,omitempty"`
	Manage Permission `json:"manage"`
	Read   Permission `json:"read"`
}

type SitesAndSensors struct {
	All    bool          `json:"all,omitempty"`
	Manage ManageSensors `json:"manage"`
	Read   Permission    `json:"read"`
}

type ManageSensors struct {
	All     bool       `json:"all,omitempty"`
	Sensors Permission `json:"sensors"`
	Sites   Permission `json:"sites"`
}

type UsersAndRoles struct {
	All    bool        `json:"all,omitempty"`
	Manage ManageUsers `json:"manage"`
	Read   Permission  `json:"read"`
}

type ManageUsers struct {
	All   bool        `json:"all,omitempty"`
	Roles UserActions `json:"roles"`
	Users UserActions `json:"users"`
}

type UserActions struct {
	All    bool       `json:"all,omitempty"`
	Create Permission `json:"create"`
	Delete Permission `json:"delete"`
	Edit   Permission `json:"edit"`
}

type User struct {
	All    bool       `json:"all,omitempty"`
	Manage ManageUser `json:"manage"`
	Read   Permission `json:"read"`
}

type ManageUser struct {
	All    bool       `json:"all,omitempty"`
	Upsert Permission `json:"upsert"`
}

type Vulnerability struct {
	All    bool       `json:"all,omitempty"`
	Manage ManageVuln `json:"manage"`
	Read   Permission `json:"read"`
}

type ManageVuln struct {
	All     bool       `json:"all,omitempty"`
	Ignore  Permission `json:"ignore"`
	Resolve Permission `json:"resolve"`
	Write   Permission `json:"write"`
}
