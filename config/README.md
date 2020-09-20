# Config

Holds configurations structs. Usually they are received by env vars.

It has general and platform specific settings as PostgreSQL, but these settings does not have dependecies from each other and can separated in the future.

Notice that we have a dependency on Viper, but we allow it because it is a "well known" package for this purpose.
