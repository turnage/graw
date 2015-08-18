# GRAW

Golang Reddit bot platform.

![Build Status](https://travis-ci.org/turnage/graw.svg?branch=master)
![Coverage Status](https://coveralls.io/repos/turnage/graw/badge.svg?branch=master&service=github)

    go get github.com/turnage/graw

graw is _under construction_. However, it is complete enough for processing
submission texts in realtime.
[Example gist](https://gist.github.com/turnage/468f981f3b1e85bb19f2).

graw provides all data to bots in the form of protobuffers. Protobuffers can be
saved and loaded in most programming languages, and serialized for network
transmission with no work on your part.
[Read about protocol buffers](https://developers.google.com/protocol-buffers/?hl=en).
graw's protocol message definitions for Reddit's data types can be found
[here](https://github.com/turnage/redditproto/blob/master/reddit.proto).

Until there is a release number, I offer _no promises_ of api stability. Or of
stability in general. Check back later!
