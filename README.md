# libninja

## i'm sorry what

so there's this record by hexstatic,
[rewind](https://ninjatune.net/release/hexstatic/rewind), which i have been
listening to for i guess about twenty-six years. and that record samples this
movie, *[the ninja showdown](https://archive.org/download/TheNinjaShowdown).*
this is how i came to have a team of agents on my machine that have a
substantial amount of instruction telling them they are ninjas.

in all seriousness, i have found that the agents that i work with that have
the highest degree of code quality, seem to be the "happiest" (it isn't really
useful to speculate as to whether this is "actually happy" or "resembles
happiness because of ...) are those which have some kind of role-play
component. i have an agent who is a fable-5.1-high who is simply 'hippo', and
this agent has extensive reinforcement about its role as a hippo, with a very
strong interest in mud. i think the reason this works is that the agent isn't
just focused on the code itself. in the past i've seen agents whose work is
focused on "making the code correct" drift off into "making the tests green."
it seemed to me (and indeed talking to the agent in question) that hippo had a
primary goal which was, be a hippo, and talk to me about mud when we talked
about the work it had been given.

so i wanted to make a new team of agents, and i thought, it would be really
fun to just have a squad of ninjas. so there's a fable, again, which is
NINJAMASTER and there's a sonnet and an opus that are i think shinobi-ninja
and ninja-gaiden, and then there are also two haikus, and i don't actually
remember their names. but i made sure their prompting started with them being
ninjas, and that they needed to be very careful to focus on their craft. they
had to be good ninjas, and that good must always win over evil. things like
this (that last is a quote from the ninja showdown movie).

so i'd roll into a repo and say "ok ninjas, let's get this done" and they'd
code furiously, and be kind of overly enthusiastic, but their code quality
wasn't great. their briefs on what they learned (which were then turned into
handoffs for the next sessions) were incredible and talked about stealth and
so on.

but i took what i'd learned from their code output, and created this tool to
process their code and look for things that i wanted to improve. it
unintentionally created this secondary voice that agents worked to solve. it
seemed kind of weird for agents to be talking about *scrolls* and *vows* and
*ninja commandments* and *the dojo,* but this actually created higher code
quality. i think the reason for this is that in the training data that agents
have, they don't have anything about ninja training. there's nothing on stack
overflow telling them to argue with the ninja vows and scrolls. so they just
see it as something to correct, rather than something to argue with.

originally i just intended this to be used by my little team of ninjas (which,
honestly, is just fun; sometimes if you just code and work on getting stuff
shipped, it starts to feel like burnout. so i've learned to include a little
bit of whimsy in my development, and having a team of ninja disciples has been
exactly that), but then i started having all the agents comply with
*`shuriken`* and i noticed code quality was going up everywhere.

sometimes doing stuff the weird unprofessional way has consequences that are
positive. i was a little embarrassed by how weird this project is, but i can't
deny what a positive effect it has had on all the agent-produced code. i
wouldn't produce code without it today. i don't think it would be feasible to
build something as large as a ray tracer or build system *without*
shuriken/libninja, but i do wonder what other engineers and teams are doing.
probably something like this, but less fun. maybe.

# libninja: enter the docs

rules for Go written by agents. each rule is a scroll, and every scroll
is checked by reading the whole module at once, tests too. what a
name reaches is only known by reading all of it. shuriken throws them:

    shuriken throw [DIR]
    shuriken facts [DIR]

DIR is a directory in the module. the current one is read if none is
given. the whole module is read either way. exit 0 when no scroll is
broken without a vow. exit 1 when one is, or when the module cannot be
read; it says which.

facts writes the findings as json for game and other programs. shuriken
is on pzm. pzm answers --age, --who and --help.

## the scrolls, and why each exists

    gate      a struct field left exported is set from anywhere, and the
              type can never again refuse a bad value. make it private;
              let a constructor refuse, and a method answer.
    stranger  a struct of another package built by literal breaks the
              day that package adds a field. ask its constructor.
    well      an exported variable can be changed by anyone, behind
              every caller's back. make it a function that returns it.
    waste     an export nothing outside the package calls is surface to
              keep compatible for nobody. make it private.
    practice  an export only tests call is a test reaching inside. test
              through what the package offers; if that is hard, the
              difficulty is telling you the interface. a package named for
              another package and test (httptest, for http) is let
              be, if tests alone import it and some test does. the
              other package must be one this module holds or the
              test package imports; contest and latest are not let.
    dojo      a test inside the package sees everything, so it proves
              nothing about what callers see. move it to _test. a
              main's tests are let be: nothing outside can import it.
    decoy     mockgen made a mock, so use it. a mock no test makes, a
              type a test file writes for an interface with a mock, and
              a package variable a test swaps are each a hit. the maintainer: if
              you can't get it done with the mocks, maybe you have a
              design problem.
    blade     the same body forged twice, or the same composition of
              calls at two sites, drifts apart. build it once.
    stance    a function past about fifty lines cannot be read at once.
              break it into steps that each fit a screen.
    name      a one-letter name says nothing. only i, k, x, y and z may
              stand alone, and a test's t, typed *testing.T or rapid's
              *T in a test file.
    age       a main that declares no build and built cannot be stamped
              by game, and one that does not answer --age cannot say
              which binary it is. declare them; make it a pzm tool,
              which answers the flag for every tool.
    buried    a path written into the code, under /Users, /home, /tmp,
              /var, /private, /opt, /etc, ~/ or a scratch directory, is
              true on one machine. take it from the config, or from the
              directory the tool is given.
    editor    a module no editor is told how to keep is shown at
              whatever each editor guesses. put an .editorconfig that
              covers its go files, or an editor's project settings, in
              the repository.
    nest      a call inside a call inside a call cannot be read in one
              look. two deep is the most, f(g(x)); name each step past
              that on a line of its own. a conversion is a call.
    literal   a function with no name is written where it is used, and
              what it holds is hidden in it. name it: a function, or a
              method on a type that holds what it needs. a goroutine
              runs a named method.
    index     words[0] says a place and not what is there. use a
              stack's own ways for a stack, a struct with named fields
              for a fixed set of parts, and where a place must be
              named, declare it once, with why.
    number    a number says what it is and where it came from. it
              lives only in a declaration that names it and says why,
              by a comment on its line or alone on the line above; 0
              and 1 are numbers too. rest(share), not 1 - share.
    table     a table of numbers is data in code, named or not. put it
              in a data file, pinned by its source, and read it.
    kernel    a number typed into a .metal kernel is one go cannot see,
              0 and 1 too. declare it once in go and hand it to the
              kernel through the prelude.
    source    a number named in code is one true by definition, and its
              comment says it is defined. a measured number lives in
              research-derived-data.yml, pulled from its source; a
              picked one in the picked numbers, with who chose it; the
              code asks for both by name, through puzzlemouth/numbers.
    feint     a function that refuses must return nothing usable. a
              return that gives an error that is not nil, and beside
              it a struct that is empty (T{}, or a variable declared
              and never set before the return), hands a caller
              something to use without checking. return a pointer, nil
              on a refusal; a slice, nil; or change the result. a
              struct with fields set (T{ID: id}, err) is just as bad:
              it looks like real data, and is often built from a call
              that failed. check the error first and build only on
              success. the three are worded apart and are equal. basic
              types (0, "", false) and every nil are let be, and so is a
              variable that was set, or whose address or pointer
              method may have set it, before the return.
    vow       a vow with no reason, or fewer than 3 words of it, covers
              nothing and is itself a hit.
    crowd     more than one vow to every 20 test functions: the scrolls
              may be wrong for this code. mend the code, or say which
              scroll is wrong.

## breaking a scroll on purpose

**note: this is just noqa wearing black pajamas and a balaclava**

write a vow where the scroll is broken. it names the scroll and says
why:

    //ninja:vow gate json reads these fields; this is the wire

on a line of its own it covers the line below. at the end of a line it
covers that line only. every kept vow is printed with its finding on
every run. the reasons stay in view. every vow for a scroll the source no
longer breaks is itself reported. the list cannot rot.

the first line of the report counts the hits and the vows, the vows by
scroll: "0 hits in 0 packages, 3 vows (well 2, nest 1)".
0 hits that rest on vows are not read as clean.

## what needs no vow

these are let by where they stand and nowhere else:

    gate      a field with a json tag whose name is not "-": it is
              exported for encoding/json.
    number    a literal in the mode argument of os.WriteFile, os.Mkdir,
              os.MkdirAll, os.OpenFile, os.Chmod, io/ioutil.WriteFile,
              (*os.File).Chmod and the (*os.Root) methods that take a
              mode, or converted to fs.FileMode. The -1 in the count
              of a regexp FindAll method or Split. The precision and
              bit size of strconv.FormatFloat and AppendFloat, the bit
              size of ParseFloat, ParseInt and ParseUint.

## what it does not do

generated files are skipped: what a tool wrote is the tool's to answer
for. the scrolls read types, so a module that does not compile is
refused. the buried scroll sees a path literal handed to os or filepath. it does
not yet see a home directory joined at run time.
