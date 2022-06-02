# TLM (Trace / Log / Metrics)

Naming things is hard and time consuming... and so is online debates over what is the best trace library, logger, and metrics system.

TLM is there to help or at least delay said discussions for some amount of time. Best thought of as a mapper library, the point is to be able to implement tracing, logging, and metrics once in your program or library once... then choose the actual implementations for those in the background. It may not win any arguments, but at least you can simply swap out the backend implementations as you want.

At best, TLM offers a consistant and light weight front end/abstraction/mapper to different libraries so as new versions arrive, security patches cause code changes, libraries get archived or created or abandoned, or curiosity on the output/functionality of another library... you don't need to update your code, just update the library.
At worst, it just repeats https://xkcd.com/927/ and adds an additional library for people to debate over.

In general: provided as-is, with a license in a different file.

## Support

### Tracers

TODO...

### Loggers

- [Logrus](https://github.com/Sirupsen/logrus)
- [ZAP](https://github.com/uber-go/zap) (Eventually)
- More??

### Metrics

TODO...
