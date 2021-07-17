<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.28"></a>
## [v1.0.28] - 2021-05-05

<a name="v1.0.27"></a>
## [v1.0.27] - 2020-10-01
### Block
- [f7587f4](https://github.com/go-jira/go-jira/commit/f7587f43f12bcf0b47e52a5abe775daf6cbb3229): reverse order of arguments
### Chore
- [2c7a9b2](https://github.com/go-jira/go-jira/commit/2c7a9b283025202428db629b1a9ecdc63a9b704f): V1.0.27 changelog bump
### Templates
- [0e3082f](https://github.com/go-jira/go-jira/commit/0e3082fab6e12a337f5fe26c3e2dec5cb51425d8): add wrap helper function

<a name="v1.0.26"></a>
## [v1.0.26] - 2020-09-11
### Chore
- [c3d22b7](https://github.com/go-jira/go-jira/commit/c3d22b765a6f3cd93445033da5c19fc0feaeaece): v1.0.26 changelog bump
### Cicd
- [96ec333](https://github.com/go-jira/go-jira/commit/96ec3339a4cc810da20450a9d9e91612c2b9aad4): automated releases fixes
- [31f7b03](https://github.com/go-jira/go-jira/commit/31f7b0397890388947f2312cf42af494c7a6979f): automated changelog and release

<a name="v1.0.25"></a>
## [v1.0.25] - 2020-09-11
### Bugfix
- [aa8dae7](https://github.com/go-jira/go-jira/commit/aa8dae7c5b7035e086bd60b3d354ffa43c30caf7): only build jira tool with gox
### Chore
- [578c44c](https://github.com/go-jira/go-jira/commit/578c44c628e3134e4d46f3250baf46d6b054cfe8): v1.0.25 release
### Fix(Changelog)
- [ff5decc](https://github.com/go-jira/go-jira/commit/ff5decc114b297e9b393f8d4af72bbace0037c73): fix changelog version
### Tests
- [a8c961f](https://github.com/go-jira/go-jira/commit/a8c961fe19f424df3fdbe108a374cc56b8ff9fe0): rework passive tests into native go tests

<a name="v1.0.24"></a>
## [v1.0.24] - 2020-09-04
### Cicd
- [d093bcf](https://github.com/go-jira/go-jira/commit/d093bcf63adbd1d4e88640307aa8a5c8669ac535): deflake tests
### Tests
- [3bc5e42](https://github.com/go-jira/go-jira/commit/3bc5e42bd0186dbc5c47f022b9528207140fa297): transition if under review
### Transition
- [3c1c4d9](https://github.com/go-jira/go-jira/commit/3c1c4d95e199a717499f1f4259649152a6832e9f): map field name to id
### Username-Deprecation
- [6a27e28](https://github.com/go-jira/go-jira/commit/6a27e28c61c45f4b2a6aff473cf28852a2df64a2): use email and display names
### Pull Requests
- Merge pull request [#367](https://github.com/go-jira/go-jira/issues/367) from bbkane/master
- Merge pull request [#349](https://github.com/go-jira/go-jira/issues/349) from aszenz/patch-1
- Merge pull request [#355](https://github.com/go-jira/go-jira/issues/355) from go-jira/vanniktech-patch-1
- Merge pull request [#323](https://github.com/go-jira/go-jira/issues/323) from tjamet/issue-comment


<a name="v1.0.23"></a>
## [v1.0.23] - 2020-02-26
### Add Sprig Template Functions, Replaces [#215] Http
- [719f7a6](https://github.com/go-jira/go-jira/commit/719f7a68a7f2c01e428a1ad3519a611c92268d27): //masterminds.github.io/sprig/
 -  [#215](https://github.com/go-jira/go-jira/issues/215)### All
- [31c113d](https://github.com/go-jira/go-jira/commit/31c113d1baf2dba814bca3e1dcc519ab8b0269e9): unindent some code
- [f125ef3](https://github.com/go-jira/go-jira/commit/f125ef3fa9c7a64e8dfda9a643cadf0241b09bc7): convert to a Go module
### CI
- [664c5ca](https://github.com/go-jira/go-jira/commit/664c5cad246cbd4c861b615eb567d3874151d1a1): test on Go 1.12.x, cleanup
### Docs
- [d8189f0](https://github.com/go-jira/go-jira/commit/d8189f0a018d1afd364237e51ca8ed43ea1aabb1): update pass documentation with password-name
### Fixes #228: Make '
- [52a577e](https://github.com/go-jira/go-jira/commit/52a577ea48afea9efb7a1f4163301129a66f7b76): ' gpg files temporary to fix go mod
 - Fixes [#228](https://github.com/go-jira/go-jira/issues/228)### Merge Branch 'Make-Password-Source-Binary-Exchangeable' Of Https
- [e26fbfc](https://github.com/go-jira/go-jira/commit/e26fbfcb142f2ce8c7c33a977d4cf0b436d743eb): //github.com/patrickpichler/jira into patrickpichler-make-password-source-binary-exchangeable
### README
- [098d963](https://github.com/go-jira/go-jira/commit/098d963881322c2b2efba48ef6a39f235bdae881): trim down the content
### T
- [d237e86](https://github.com/go-jira/go-jira/commit/d237e86cda3812b23f432e90d120ec21e749a854): rename to _t to fix module support
 - Fixes [#228](https://github.com/go-jira/go-jira/issues/228)### Update All Usage Of User.Name To User.AccountId For Privacy Migration: Https
- [a26683e](https://github.com/go-jira/go-jira/commit/a26683e01dc7e161e735b1b387d1633bc32da2fe): //developer.atlassian.com/cloud/jira/platform/deprecation-notice-user-privacy-api-migration-guide/
### Pull Requests
- Merge pull request [#318](https://github.com/go-jira/go-jira/issues/318) from jrschumacher/patch-1
- Merge pull request [#317](https://github.com/go-jira/go-jira/issues/317) from go-jira/privacy-migration
- Merge pull request [#302](https://github.com/go-jira/go-jira/issues/302) from go-jira/simplify-template-tables
- Merge pull request [#292](https://github.com/go-jira/go-jira/issues/292) from pdecat/cache_password
- Merge pull request [#301](https://github.com/go-jira/go-jira/issues/301) from go-jira/allow-issue-ints
- Merge pull request [#286](https://github.com/go-jira/go-jira/issues/286) from patrickpichler/add-gopass-instructions-to-readme
- Merge pull request [#285](https://github.com/go-jira/go-jira/issues/285) from patrickpichler/add-gopass-support
- Merge pull request [#283](https://github.com/go-jira/go-jira/issues/283) from go-jira/sprig
- Merge pull request [#273](https://github.com/go-jira/go-jira/issues/273) from acaloiaro/master
- Merge pull request [#282](https://github.com/go-jira/go-jira/issues/282) from pcockwell/fix/choose-direct-transition-match-if-available
- Merge pull request [#280](https://github.com/go-jira/go-jira/issues/280) from go-jira/cut_v1_0_21
- Merge pull request [#275](https://github.com/go-jira/go-jira/issues/275) from go-jira/remove_gopkg_in
- Merge pull request [#278](https://github.com/go-jira/go-jira/issues/278) from go-jira/update-figtree
- Merge pull request [#276](https://github.com/go-jira/go-jira/issues/276) from go-jira/fix_228
- Merge pull request [#266](https://github.com/go-jira/go-jira/issues/266) from mbethke/fix-multiline-worklog-comment
- Merge pull request [#263](https://github.com/go-jira/go-jira/issues/263) from kojustin/master
- Merge pull request [#253](https://github.com/go-jira/go-jira/issues/253) from mvdan/module
- Merge pull request [#240](https://github.com/go-jira/go-jira/issues/240) from jgraglia/patch-1
- Merge pull request [#219](https://github.com/go-jira/go-jira/issues/219) from kerhac/master
- Merge pull request [#236](https://github.com/go-jira/go-jira/issues/236) from CodeLingoBot/rewrite
- Merge pull request [#245](https://github.com/go-jira/go-jira/issues/245) from justmiles/211


<a name="v1.0.22"></a>
## [v1.0.22] - 2019-09-30
### All
- [bb9790f](https://github.com/go-jira/go-jira/commit/bb9790f28783c1b82a3685a9c4657241e906826a): unindent some code
- [89fe2ec](https://github.com/go-jira/go-jira/commit/89fe2ecf16709511c3e04e02f7c7906a5ac6865a): convert to a Go module
### CI
- [80743e4](https://github.com/go-jira/go-jira/commit/80743e4da870a1febcc65d18d08242bb201b744d): test on Go 1.12.x, cleanup
### Docs
- [48c15e2](https://github.com/go-jira/go-jira/commit/48c15e2daa7b3f4c84526bd9f030828f378edfc2): update pass documentation with password-name
### Fixes #228: Make '
- [3984d0d](https://github.com/go-jira/go-jira/commit/3984d0d4848fdfe790f46ec18bd3b71782e36c32): ' gpg files temporary to fix go mod
 - Fixes [#228](https://github.com/go-jira/go-jira/issues/228)### README
- [dc9a9de](https://github.com/go-jira/go-jira/commit/dc9a9de165859057e4596aa47961e84de34b0b4b): trim down the content
### T
- [8994b42](https://github.com/go-jira/go-jira/commit/8994b42f714f8fc5b224bda8b5835f003d96ef02): rename to _t to fix module support
 - Fixes [#228](https://github.com/go-jira/go-jira/issues/228)
<a name="v1.0.21"></a>
## [v1.0.21] - 2019-09-16
### All
- [31c113d](https://github.com/go-jira/go-jira/commit/31c113d1baf2dba814bca3e1dcc519ab8b0269e9): unindent some code
- [f125ef3](https://github.com/go-jira/go-jira/commit/f125ef3fa9c7a64e8dfda9a643cadf0241b09bc7): convert to a Go module
### CI
- [664c5ca](https://github.com/go-jira/go-jira/commit/664c5cad246cbd4c861b615eb567d3874151d1a1): test on Go 1.12.x, cleanup
### Docs
- [d8189f0](https://github.com/go-jira/go-jira/commit/d8189f0a018d1afd364237e51ca8ed43ea1aabb1): update pass documentation with password-name
### Fixes #228: Make '
- [52a577e](https://github.com/go-jira/go-jira/commit/52a577ea48afea9efb7a1f4163301129a66f7b76): ' gpg files temporary to fix go mod
 - Fixes [#228](https://github.com/go-jira/go-jira/issues/228)### README
- [098d963](https://github.com/go-jira/go-jira/commit/098d963881322c2b2efba48ef6a39f235bdae881): trim down the content
### T
- [d237e86](https://github.com/go-jira/go-jira/commit/d237e86cda3812b23f432e90d120ec21e749a854): rename to _t to fix module support
 - Fixes [#228](https://github.com/go-jira/go-jira/issues/228)### Pull Requests
- Merge pull request [#275](https://github.com/go-jira/go-jira/issues/275) from go-jira/remove_gopkg_in
- Merge pull request [#278](https://github.com/go-jira/go-jira/issues/278) from go-jira/update-figtree
- Merge pull request [#276](https://github.com/go-jira/go-jira/issues/276) from go-jira/fix_228
- Merge pull request [#266](https://github.com/go-jira/go-jira/issues/266) from mbethke/fix-multiline-worklog-comment
- Merge pull request [#263](https://github.com/go-jira/go-jira/issues/263) from kojustin/master
- Merge pull request [#253](https://github.com/go-jira/go-jira/issues/253) from mvdan/module
- Merge pull request [#240](https://github.com/go-jira/go-jira/issues/240) from jgraglia/patch-1
- Merge pull request [#219](https://github.com/go-jira/go-jira/issues/219) from kerhac/master
- Merge pull request [#236](https://github.com/go-jira/go-jira/issues/236) from CodeLingoBot/rewrite
- Merge pull request [#245](https://github.com/go-jira/go-jira/issues/245) from justmiles/211
- Merge pull request [#220](https://github.com/go-jira/go-jira/issues/220) from ejsuncy/master


<a name="v1.0.20"></a>
## [v1.0.20] - 2018-08-04

<a name="v1.0.19"></a>
## [v1.0.19] - 2018-08-02
### Pull Requests
- Merge pull request [#197](https://github.com/go-jira/go-jira/issues/197) from kojiromike/spellcheck


<a name="v1.0.18"></a>
## [v1.0.18] - 2018-07-29
### They Broke Golang.Org/X/Net: ERROR: Vendor/Golang.Org/X/Net/Proxy/Socks5.Go:11:2
- [7191c77](https://github.com/go-jira/go-jira/commit/7191c7751b2d18d7f951d089fa3235acf5748d4b): use of internal package not allowed
### Pull Requests
- Merge pull request [#178](https://github.com/go-jira/go-jira/issues/178) from vergenzt/patch-1


<a name="v1.0.17"></a>
## [v1.0.17] - 2018-04-15
### [#157] Add `Password-Directory
- [06b26c9](https://github.com/go-jira/go-jira/commit/06b26c9e00384318ec7a51fa1c5ff5de63ea686b): path` to allow overriding PASSWORD_STORE_DIR from configs
 -  [#157](https://github.com/go-jira/go-jira/issues/157)### Pull Requests
- Merge pull request [#161](https://github.com/go-jira/go-jira/issues/161) from vanniktech/patch-1


<a name="v1.0.16"></a>
## [v1.0.16] - 2018-04-01
### Pull Requests
- Merge pull request [#150](https://github.com/go-jira/go-jira/issues/150) from catskul/parameterized-go-makefile
- Merge pull request [#153](https://github.com/go-jira/go-jira/issues/153) from catskul/document-shell-completion
- Merge pull request [#152](https://github.com/go-jira/go-jira/issues/152) from catskul/fix-missing-priority


<a name="v1.0.15"></a>
## [v1.0.15] - 2018-03-08
### Pull Requests
- Merge pull request [#151](https://github.com/go-jira/go-jira/issues/151) from catskul/build-instructions
- Merge pull request [#142](https://github.com/go-jira/go-jira/issues/142) from anthonyrisinger/patch-1


<a name="v1.0.14"></a>
## [v1.0.14] - 2017-11-04
### Pull Requests
- Merge pull request [#130](https://github.com/go-jira/go-jira/issues/130) from onionjake/master


<a name="v1.0.13"></a>
## [v1.0.13] - 2017-10-28
### Pull Requests
- Merge pull request [#126](https://github.com/go-jira/go-jira/issues/126) from schorsch3000/master
- Merge pull request [#129](https://github.com/go-jira/go-jira/issues/129) from blachniet/logout-help-typo-fix
- Merge pull request [#124](https://github.com/go-jira/go-jira/issues/124) from gvol/master
- Merge pull request [#128](https://github.com/go-jira/go-jira/issues/128) from mivok/escape-issuetype


<a name="v1.0.12"></a>
## [v1.0.12] - 2017-10-04

<a name="v1.0.11"></a>
## [v1.0.11] - 2017-09-26

<a name="v1.0.10"></a>
## [v1.0.10] - 2017-09-18

<a name="v1.0.9"></a>
## [v1.0.9] - 2017-09-17

<a name="v1.0.8"></a>
## [v1.0.8] - 2017-09-17

<a name="v1.0.7"></a>
## [v1.0.7] - 2017-09-15

<a name="v1.0.6"></a>
## [v1.0.6] - 2017-09-13

<a name="v1.0.5"></a>
## [v1.0.5] - 2017-09-11

<a name="v1.0.4"></a>
## [v1.0.4] - 2017-09-08

<a name="v1.0.3"></a>
## [v1.0.3] - 2017-09-06

<a name="v1.0.2"></a>
## [v1.0.2] - 2017-09-06

<a name="v1.0.1"></a>
## [v1.0.1] - 2017-09-06

<a name="v1.0.0"></a>
## [v1.0.0] - 2017-09-05

<a name="v0.1.15"></a>
## [v0.1.15] - 2017-08-25
### Pull Requests
- Merge pull request [#104](https://github.com/go-jira/go-jira/issues/104) from wrouesnel/keyring-update
- Merge pull request [#90](https://github.com/go-jira/go-jira/issues/90) from bbaugher/master


<a name="v0.1.14"></a>
## [v0.1.14] - 2017-05-10

<a name="v0.1.13"></a>
## [v0.1.13] - 2017-04-24
### Pull Requests
- Merge pull request [#78](https://github.com/go-jira/go-jira/issues/78) from davidreuss/generic-issuelink
- Merge pull request [#77](https://github.com/go-jira/go-jira/issues/77) from davidreuss/fix-start-parameter-for-pagination


<a name="v0.1.12"></a>
## [v0.1.12] - 2017-03-22
### Pull Requests
- Merge pull request [#74](https://github.com/go-jira/go-jira/issues/74) from clausb/BrowseOnWindows


<a name="v0.1.11"></a>
## [v0.1.11] - 2017-02-26

<a name="v0.1.10"></a>
## [v0.1.10] - 2017-02-08
### Doc Tweak
- [e6faa4e](https://github.com/go-jira/go-jira/commit/e6faa4eab1a8d6a7fb624b79bb58a641d02e876b): add info for setting username
### Merge Branch 'Master' Of Github.Com
- [63bc2ae](https://github.com/go-jira/go-jira/commit/63bc2ae15a2ebafa16861965951800e0d5c122bd): Netflix-Skunkworks/go-jira
### Refactor Password Source, Allow For "Pass" To Be Used, Update Tests To Use `Password-Source
- [cb70941](https://github.com/go-jira/go-jira/commit/cb70941aad2b8198f5c8ad8d1e1a7a98dc820cd9): pass`
### Pull Requests
- Merge pull request [#65](https://github.com/go-jira/go-jira/issues/65) from mlbright/patch-1
- Merge pull request [#64](https://github.com/go-jira/go-jira/issues/64) from astrostl/patch-2
- Merge pull request [#62](https://github.com/go-jira/go-jira/issues/62) from astrostl/patch-1


<a name="v0.1.9"></a>
## [v0.1.9] - 2016-12-18
### Fix(Http)
- [b326623](https://github.com/go-jira/go-jira/commit/b326623ed22677a3ff76d2c4c67bb7ca7ecb3877): Add proxy transport
- [72c78c6](https://github.com/go-jira/go-jira/commit/72c78c6c1c63a70d837c8e367754792c8a30ae06): Add proxy transport
### Merge Branch 'Master' Of Github.Com
- [ac515e2](https://github.com/go-jira/go-jira/commit/ac515e2743e1bcf5f492a0e25d2b084f3311f0d0): Netflix-Skunkworks/go-jira
### Pull Requests
- Merge pull request [#61](https://github.com/go-jira/go-jira/issues/61) from sylus/feature-proxy
- Merge pull request [#60](https://github.com/go-jira/go-jira/issues/60) from facundoolano/patch-1


<a name="v0.1.8"></a>
## [v0.1.8] - 2016-11-24
### Pull Requests
- Merge pull request [#53](https://github.com/go-jira/go-jira/issues/53) from jshirley/master


<a name="v0.1.7"></a>
## [v0.1.7] - 2016-08-24
### Pull Requests
- Merge pull request [#52](https://github.com/go-jira/go-jira/issues/52) from dbrower/master


<a name="v0.1.6"></a>
## [v0.1.6] - 2016-08-21

<a name="v0.1.5"></a>
## [v0.1.5] - 2016-08-21

<a name="v0.1.4"></a>
## [v0.1.4] - 2016-08-12

<a name="v0.1.3"></a>
## [v0.1.3] - 2016-07-30
### Pull Requests
- Merge pull request [#24](https://github.com/go-jira/go-jira/issues/24) from mikepea/edit_template_common


<a name="v0.1.2"></a>
## [v0.1.2] - 2016-06-29

<a name="v0.1.1"></a>
## [v0.1.1] - 2016-06-28
### Merge Branch 'Master' Of Github.Com
- [dd0f5ef](https://github.com/go-jira/go-jira/commit/dd0f5efd3247157686c2c88817d3ad375de399ea): Netflix-Skunkworks/go-jira
### Pull Requests
- Merge pull request [#39](https://github.com/go-jira/go-jira/issues/39) from mikepea/system_template_dir
- Merge pull request [#38](https://github.com/go-jira/go-jira/issues/38) from jglick/patch-1
- Merge pull request [#37](https://github.com/go-jira/go-jira/issues/37) from tobyjoe/add-resource-expansion
- Merge pull request [#35](https://github.com/go-jira/go-jira/issues/35) from QuinnyPig/fix-readme
- Merge pull request [#34](https://github.com/go-jira/go-jira/issues/34) from jonathanio/fix/issuetypes-url-escaping


<a name="v0.1.0"></a>
## [v0.1.0] - 2016-01-29
### Add Component/Components Support
- [595a521](https://github.com/go-jira/go-jira/commit/595a5212b43be28e01a0d5a6cf98a8de89383e41): add and list for now.
### Merge Branch 'Master' Of Github.Com
- [9e90376](https://github.com/go-jira/go-jira/commit/9e90376816c295d3a75b4f51703c24fd95873625): Netflix-Skunkworks/go-jira
- [382bf4f](https://github.com/go-jira/go-jira/commit/382bf4faeb17198b54950a05b0fdb3e126af8d73): Netflix-Skunkworks/go-jira
### Pull Requests
- Merge pull request [#33](https://github.com/go-jira/go-jira/issues/33) from mikepea/make_jirad
- Merge pull request [#31](https://github.com/go-jira/go-jira/issues/31) from mikepea/component_mgmt
- Merge pull request [#30](https://github.com/go-jira/go-jira/issues/30) from mikepea/unwatch_support
- Merge pull request [#26](https://github.com/go-jira/go-jira/issues/26) from mikepea/vote_support


<a name="v0.0.20"></a>
## [v0.0.20] - 2016-01-21
### Correct Naming Of Parameter
- [8e66246](https://github.com/go-jira/go-jira/commit/8e662462dac6cccdf8af277797777caeff3ad2bc): set/add/remove are actions.
### Pull Requests
- Merge pull request [#27](https://github.com/go-jira/go-jira/issues/27) from blalor/insecure-skip-verify
- Merge pull request [#21](https://github.com/go-jira/go-jira/issues/21) from mikepea/label_command
- Merge pull request [#22](https://github.com/go-jira/go-jira/issues/22) from mikepea/library_break_out
- Merge pull request [#20](https://github.com/go-jira/go-jira/issues/20) from mikepea/add_join_template_func


<a name="0.0.19"></a>
## [0.0.19] - 2015-12-09

<a name="0.0.18"></a>
## [0.0.18] - 2015-12-03

<a name="0.0.17"></a>
## [0.0.17] - 2015-12-03

<a name="0.0.16"></a>
## [0.0.16] - 2015-11-23

<a name="0.0.15"></a>
## [0.0.15] - 2015-11-23

<a name="0.0.14"></a>
## [0.0.14] - 2015-11-17
### Pull Requests
- Merge pull request [#16](https://github.com/go-jira/go-jira/issues/16) from oschrenk/fix-typo
- Merge pull request [#14](https://github.com/go-jira/go-jira/issues/14) from mikepea/ls_with_updated


<a name="0.0.13"></a>
## [0.0.13] - 2015-09-19

<a name="0.0.12"></a>
## [0.0.12] - 2015-09-18

<a name="0.0.11"></a>
## [0.0.11] - 2015-09-16

<a name="0.0.10"></a>
## [0.0.10] - 2015-09-15

<a name="0.0.9"></a>
## [0.0.9] - 2015-09-15
### Allow "Abort
- [80b6f5a](https://github.com/go-jira/go-jira/commit/80b6f5a198fbe17aa0245c470d47c2988d8624c3): true" to be set while editing to cancel the edit operation

<a name="0.0.8"></a>
## [0.0.8] - 2015-07-31
### Pull Requests
- Merge pull request [#11](https://github.com/go-jira/go-jira/issues/11) from mikepea/max_results_option

### Note

that testing against our JIRA instance, in a project with
more than 1000 open issues, suggests that the JIRA has an internal
limit of 1000 results in a single query.


<a name="0.0.7"></a>
## [0.0.7] - 2015-07-01
### Merge Branch 'Master' Of Github.Com
- [b72040b](https://github.com/go-jira/go-jira/commit/b72040bfd413bcdc88ca2c3f6843a7f6dee2e898): Netflix-Skunkworks/go-jira
### Pull Requests
- Merge pull request [#9](https://github.com/go-jira/go-jira/issues/9) from nelfin/quickfix/take-user


<a name="0.0.6"></a>
## [0.0.6] - 2015-02-27
### Set JIRA_OPERATION To "View" When No Operation Used (Ie
- [8040746](https://github.com/go-jira/go-jira/commit/8040746bcf6aac6e3ff6e419c349661a8fc9bf99): jira GOJIRA-123)

<a name="0.0.5"></a>
## [0.0.5] - 2015-02-21

<a name="0.0.4"></a>
## [0.0.4] - 2015-02-19

<a name="0.0.3"></a>
## [0.0.3] - 2015-02-19
### [Issue #8] Detect X-Seraph-Loginreason
- [f3feff7](https://github.com/go-jira/go-jira/commit/f3feff796fbecca6477ecbc1e9dae6a2b78e755c): AUTHENTICATION_DENIED header to catch login failures
 -  [#8](https://github.com/go-jira/go-jira/issues/8)### Pull Requests
- Merge pull request [#7](https://github.com/go-jira/go-jira/issues/7) from jaybuff/empty-projects


<a name="0.0.2"></a>
## [0.0.2] - 2015-02-18

<a name="0.0.1"></a>
## 0.0.1 - 2015-02-18
### Adding Commands
- [18f10fd](https://github.com/go-jira/go-jira/commit/18f10fd12521584c7d85b20dff2b1c2da0854cb9): * create * dups * blocks * watch
### Merge Branch 'Nil-Assignee' Of Https
- [25539ef](https://github.com/go-jira/go-jira/commit/25539efedda0e06e103fc942e16405c0c09ba621): //github.com/jaybuff/go-jira into jaybuff-nil-assignee
### Work In Progress, Minor Refactor.  Added Commands
- [acbc24b](https://github.com/go-jira/go-jira/commit/acbc24b2096f31a5805fa48984724a4a6c1da431): * login * editmeta ISSUE * edit ISSUE * issuetypes [-p PROJECT] * createmeta [-p PROJECT] [-i ISSUETYPE] * transitions ISSUE
### Pull Requests
- Merge pull request [#2](https://github.com/go-jira/go-jira/issues/2) from jaybuff/clean-up


[Unreleased]: https://github.com/go-jira/go-jira/compare/v1.0.28...HEAD
[v1.0.28]: https://github.com/go-jira/go-jira/compare/v1.0.27...v1.0.28
[v1.0.27]: https://github.com/go-jira/go-jira/compare/v1.0.26...v1.0.27
[v1.0.26]: https://github.com/go-jira/go-jira/compare/v1.0.25...v1.0.26
[v1.0.25]: https://github.com/go-jira/go-jira/compare/v1.0.24...v1.0.25
[v1.0.24]: https://github.com/go-jira/go-jira/compare/v1.0.23...v1.0.24
[v1.0.23]: https://github.com/go-jira/go-jira/compare/v1.0.22...v1.0.23
[v1.0.22]: https://github.com/go-jira/go-jira/compare/v1.0.21...v1.0.22
[v1.0.21]: https://github.com/go-jira/go-jira/compare/v1.0.20...v1.0.21
[v1.0.20]: https://github.com/go-jira/go-jira/compare/v1.0.19...v1.0.20
[v1.0.19]: https://github.com/go-jira/go-jira/compare/v1.0.18...v1.0.19
[v1.0.18]: https://github.com/go-jira/go-jira/compare/v1.0.17...v1.0.18
[v1.0.17]: https://github.com/go-jira/go-jira/compare/v1.0.16...v1.0.17
[v1.0.16]: https://github.com/go-jira/go-jira/compare/v1.0.15...v1.0.16
[v1.0.15]: https://github.com/go-jira/go-jira/compare/v1.0.14...v1.0.15
[v1.0.14]: https://github.com/go-jira/go-jira/compare/v1.0.13...v1.0.14
[v1.0.13]: https://github.com/go-jira/go-jira/compare/v1.0.12...v1.0.13
[v1.0.12]: https://github.com/go-jira/go-jira/compare/v1.0.11...v1.0.12
[v1.0.11]: https://github.com/go-jira/go-jira/compare/v1.0.10...v1.0.11
[v1.0.10]: https://github.com/go-jira/go-jira/compare/v1.0.9...v1.0.10
[v1.0.9]: https://github.com/go-jira/go-jira/compare/v1.0.8...v1.0.9
[v1.0.8]: https://github.com/go-jira/go-jira/compare/v1.0.7...v1.0.8
[v1.0.7]: https://github.com/go-jira/go-jira/compare/v1.0.6...v1.0.7
[v1.0.6]: https://github.com/go-jira/go-jira/compare/v1.0.5...v1.0.6
[v1.0.5]: https://github.com/go-jira/go-jira/compare/v1.0.4...v1.0.5
[v1.0.4]: https://github.com/go-jira/go-jira/compare/v1.0.3...v1.0.4
[v1.0.3]: https://github.com/go-jira/go-jira/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/go-jira/go-jira/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/go-jira/go-jira/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/go-jira/go-jira/compare/v0.1.15...v1.0.0
[v0.1.15]: https://github.com/go-jira/go-jira/compare/v0.1.14...v0.1.15
[v0.1.14]: https://github.com/go-jira/go-jira/compare/v0.1.13...v0.1.14
[v0.1.13]: https://github.com/go-jira/go-jira/compare/v0.1.12...v0.1.13
[v0.1.12]: https://github.com/go-jira/go-jira/compare/v0.1.11...v0.1.12
[v0.1.11]: https://github.com/go-jira/go-jira/compare/v0.1.10...v0.1.11
[v0.1.10]: https://github.com/go-jira/go-jira/compare/v0.1.9...v0.1.10
[v0.1.9]: https://github.com/go-jira/go-jira/compare/v0.1.8...v0.1.9
[v0.1.8]: https://github.com/go-jira/go-jira/compare/v0.1.7...v0.1.8
[v0.1.7]: https://github.com/go-jira/go-jira/compare/v0.1.6...v0.1.7
[v0.1.6]: https://github.com/go-jira/go-jira/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/go-jira/go-jira/compare/v0.1.4...v0.1.5
[v0.1.4]: https://github.com/go-jira/go-jira/compare/v0.1.3...v0.1.4
[v0.1.3]: https://github.com/go-jira/go-jira/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/go-jira/go-jira/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/go-jira/go-jira/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/go-jira/go-jira/compare/v0.0.20...v0.1.0
[v0.0.20]: https://github.com/go-jira/go-jira/compare/0.0.19...v0.0.20
[0.0.19]: https://github.com/go-jira/go-jira/compare/0.0.18...0.0.19
[0.0.18]: https://github.com/go-jira/go-jira/compare/0.0.17...0.0.18
[0.0.17]: https://github.com/go-jira/go-jira/compare/0.0.16...0.0.17
[0.0.16]: https://github.com/go-jira/go-jira/compare/0.0.15...0.0.16
[0.0.15]: https://github.com/go-jira/go-jira/compare/0.0.14...0.0.15
[0.0.14]: https://github.com/go-jira/go-jira/compare/0.0.13...0.0.14
[0.0.13]: https://github.com/go-jira/go-jira/compare/0.0.12...0.0.13
[0.0.12]: https://github.com/go-jira/go-jira/compare/0.0.11...0.0.12
[0.0.11]: https://github.com/go-jira/go-jira/compare/0.0.10...0.0.11
[0.0.10]: https://github.com/go-jira/go-jira/compare/0.0.9...0.0.10
[0.0.9]: https://github.com/go-jira/go-jira/compare/0.0.8...0.0.9
[0.0.8]: https://github.com/go-jira/go-jira/compare/0.0.7...0.0.8
[0.0.7]: https://github.com/go-jira/go-jira/compare/0.0.6...0.0.7
[0.0.6]: https://github.com/go-jira/go-jira/compare/0.0.5...0.0.6
[0.0.5]: https://github.com/go-jira/go-jira/compare/0.0.4...0.0.5
[0.0.4]: https://github.com/go-jira/go-jira/compare/0.0.3...0.0.4
[0.0.3]: https://github.com/go-jira/go-jira/compare/0.0.2...0.0.3
[0.0.2]: https://github.com/go-jira/go-jira/compare/0.0.1...0.0.2
