# Administration

* POST `/admin/staff/new`
* GET `/admin/staff/roster?page=<number>` All staff

* GET `/admin/staff/profile/{name}` Show a staff's profile
* PUT `/admin/staff/profile/{name}` Restore a deleted staff
* PATCH `/admin/staff/profile/{name}` Update staff's profile
* DELETE `/admin/staff/profile/{name}?rmvip=true|false` Delete a staff

* GET `/admin/vip` Show all myft accounts that are granted VIP.
* PUT `/admin/vip/{myftId}` Grant vip to a myft account
* DELETE `/admin/vip/{myftId}` Delete vip status of a myft account