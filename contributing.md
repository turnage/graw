# Contributing

graw welcomes community contributions!

graw is a load bearing package. Many users depend on graw, and graw promises not to break
users. The most important thing to keep in mind when changing graw is that we
*cannot* break users.

## Never Break Users

Any source code that builds at any version of graw after 1.0 must continue to build in all
later versions of graw. Any semantics that users have come to depend on, including, to some
extent, unintended semantics, must be preserved.

## Add and Maintain Automated Tests

If you add a feature to graw, great! Users may come to depend on it, and from then on we
need to ensure we do not break it. This is why I ask that all new features also come with
automated tests. Future contributors will need to ensure these still pass in order to check
in code.

## Ruthlessly Segment Commit History

Pull requests are best organized into small commits containing only related changes. For
example, if you want to add an API call and want to refactor some things to make that easier,
I would like to see independent commits for each refactor, followed the commit introducing
new behavior and its tests.

I consider most commits over 100 lines in diff unreviewable. Exceptions are self contained
new implementations and easy to review tooling-generated changes such as project wide renames.

For more reading on this philosophy, see [Why small CLs?](https://google.github.io/eng-practices/review/developer/small-cls.html)

## Get Credit

Help yourself to [contributors.md](contributors.md).
