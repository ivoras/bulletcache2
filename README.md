# Bullet Cache 2 (BC2)

A high performance, multithreaded in-memory cache database optimized for group record expiry and supporting multiple protocols (memcached, redis, gRPC, REST).

## Plan

The overarching idea is to prioritise performance, even to the point of sacrificing some features and comfort for it. In its basic functionalities, it's more similar to Memcached than Redis, with the added feature of having optional tags attached to records, which are mainly intended to *help with cache invalidation*. It will support multiple protocols accessing the same memory cache database, including Memcached, Redis, gRPC and plain HTTP.

In addition to setting and deleting records by addressing them with keys, they can be deleted in bulk by addressing them with tags, enabling bulk cache invalidation. In spirit, this is a continuation (and closure) of the project I did in my PhD thesis a long long time ago, which started as `mdcached` - the multi-domain cache server, and was developed into [Bullet Cache](https://mdcached.sourceforge.net/). I've learned a lot since then and this project should be much more robust.

The plan is to create the following:

- [ ] In-memory key-value database, with binary keys and values, with keys restricted to up to 250 bytes
- [ ] Records can have uint32 tags, up to 8 tags per record
- [ ] Records can be set in an atomic operation which sets the key, the value and the tags
- [ ] Records can be queried by keys
- [ ] Records can be queried by tags
- [ ] Records can be deleted by keys
- [ ] Records can be deleted by tags
- [ ] A memcached-compatible network protocol for easier adoption
- [ ] Implementaion of every single Memcached op from https://github.com/memcached/memcached/wiki/Commands
- [ ] Redis-compatible protocol for easier adoption, at least for simpler commands
- [ ] A gRPC-based protocol for added performance
- [ ] A HTTP REST-like network protocol for convenience, using https://github.com/valyala/fasthttp (yes, this means you'll be able to SET some data into the in-memory cache using the Memcached or Redis protocol, and serve them from memory via HTTP)
- [ ] Redis-like persistance and loading for "permanent" data

BC2 aims to be a good microservice player, configurable and usable in a container. See the `.env.template` file for configurable variables.

# Data model

AS far as users of this service are concerned, each record stored in BC2 is immutable and has the following basic structure:

- `key` - A binary array of up to 250 bytes, unique
- `value` - A binary array of arbitrary size
- `tags` - A list of up to 8 32-bit unsigned integers

Tags and keys are internally indexed to support fast operations such as these:

- Get, set and delete based on the key. Set operations usually replace existing records with the same key.
- Get and delete based on an integer tag. This enables atomic expiry of a whole group of records, as well as retriving a group of records in protocols which support it.

Tags are arbitrary integers assigned to records, with the intention that application use them to group records in some meaningful way - mostly in groups of similar data which need to be expired at the same time, e.g. related to a particular user, etc.