# GoServerTemplate

A high performance, multithreaded memory cache server with domain tagging.

Plan:

The overarching idea is to prioritise performance, even to the point of sacrificing some features and comfort for it. In its basic functionalities, it's more similar to Memcached than Redis, with the added feature of having optional tags attached to records, which are intended to help with cache invalidation. In addition to setting and deleting records by addressing them with their keys, they can be deleted in bulk by addressing them with tags. In spirit, this is a continuation (and closure) of the project I did in my PhD thesis a long long time ago, which started as `mdcached` - the multi-domain cache server, and was developed into [Bullet Cache](https://mdcached.sourceforge.net/). I've learned a lot since then and this project should be much more robust.

The plan is to create the following:

- [ ] In-memory key-value database, with binary keys and values, but keys are restricted to up to 255 bytes
- [ ] Records can have uint32 tags, up to 4 tags per record
- [ ] Records can be set in an atomic operation which sets the key, the value and the tags
- [ ] Records can be queried by keys
- [ ] Records can be queried by tags
- [ ] Records can be deleted by keys
- [ ] Records can be deleted by tags
