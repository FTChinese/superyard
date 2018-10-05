# Customer Service

## Search

* GET `/search/user?k=<name|email>&v=:value`
* GET `/search/orders?{start=YYYY-MM-DD&end=YYYY-MM-DD}` Show all orders within the specified time range

## FTC User

* GET `/ftc-user/profile/{userId}` Show a user's profile, vip status, membership
* GET `/ftc-user/profile/{userId}/orders` Show a user's orders
* GET `/ftc-user/profile/{userId}/login` Show a user's login history
