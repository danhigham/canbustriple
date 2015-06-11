package main

import (
	"github.com/danhigham/gocui"
  "strings"
  "fmt"
)

const logo string = `





                                            ''''
                                            +hhhs-
           ':+ohho/       -yyo/- .:yhho:     yhhhs+ooohhyoo:.
        '/shhhhhhhh/     .yhhhhy.  shhhh+   'yhhhhhhhhhhhhhhhs'
       -hhhhs+::sho    .ohhhhy/'  .yhhhhh   /hhhsooo::'''-hhhh/
     'ohhhh:    ':'  .+hhhhsyos/ 'shhhhhh  'yhho /hhhh- .ohhhs' /+/'  '+/:   './/++-
    .yhhhy.        .ohhhho.ohhh: ohhhhhhh  +hhs -yhhho/yhhho:  .yhh-  ohhy' /yhhsyhh
   'shhhy.     '  /hhhho. :hhho /hhho+hhh/-hhy''hhhh+:syhhy:  .yhh-  ohhy. ohhh+-/o+ooooooooooo:
   /hhhh/    -++:shhhy:++yyhhh.:hhh+ /hhhsyhh-'shhhy'  '/hhh/ shh: -shhy.  ohhhhhyo:::::/oohhhhh-
   shhhhsooohh+:hhhho/hhhhhhhy.yhhs' 'yhhhhh+'+hhhhh+:::ohhhs+hhy:shhhh- :yho'-:hhh-        ./y:
   /hhhhhhhh+' ohhho '..-yhhh::hhh-   +hhhhs .yhhyyhhhhhhhh+-shhhhsohhhyohhhhssyhy/'
    .+yy+/-     :yy.    -hhhh: :o+     -++/'  -+o/'/+sys+:.  .+++- '/s+: .:+++++-'
                        -hhhh:                                      '.'
 ..''''''-----/+++++++- .yhhhs         ''                           ohhy+.
'yhhhhhhhhhhhhhhhhhhhhhy..:ssss. '-/+sshhs/'                       -hhhhhy'
:hhhhhhhhhhhhhhhhhhhhhhho     ':+yhhhhhhhhh-                      'yhhhhh:
'yhhhhhhhhhhhhhhysssssss-    -shhhhhhhhhhho'                     'shhhhho                 '''
 ''::::/hhhhhhh+    .-'    :shhhs+::''''''  --'                  ohhhhho'               -yhhhs+'
      'shhhhhh/     yhhy-'ohho-'-::.'      -hhhs-    '':.'      +hhhhh:    '':::::.'     ohhhhhs'
      +hhhhhh+     -yhhhyyho.  :hhhhh-    'shhhhh../shhhhho.   +hhhhh/   :+yhhhhhhhy.   'shhhhhh-
     /hhhhhh+     'shhhhhh/   -yhhhhy     +hhhhhsohhhhhhhhh+  -hhhhh/  -shhhhhhhhhhh+   -hhhhhhs'
    -hhhhhhs      ohhhhhy-   'hhhhhy.    /hhhhhhhys/ohhhhhh+ -hhhhho''shhhhhy/+yhhhh+   shhhhhh:
   .hhhhhhy'     +hhhhhy.   'shhhhy.    :hhhhhhs-'  :hhhhhh-.yhhhho'-yhhhhy:  'shhhh'  +hhhhhho
  .yhhhhhy-     :yhhhhh-   .yhhhhy.    :hhhhhh+'   .hhhhhh+'shhhhy -hhhhh+/+/+yhhhs' 'ohhhhhho
  ohhhhhh/     -hhhhhy:    yhhhhy.    .yhhhhhh-   :yhhhhho'/hhhhh:'yhhhhs:hhhhhhy:  :shhhhhh:
 +hhhhhhs     .yhhhhy:    +hhhhh-   -+hhhhhhhh/':ohhhhhh/' yhhhhy :hhhhh: -+++-' ':yhhhhhho.
.hhhhhhy'    'shhhhh/    'hhhhhy::oyhhhhhhhhhhhhhhhhhhh-   yhhhhy :hhhhho-   .-+shhhhhhho.
ohhhhhh-     -hhhhho     'shhhhhhhhyhhhhhhs/hhhhhhhhh+'    yhhhhy 'shhhhhhhhhhhhhhhhhy+.
:hhhhh:      'oyhhy'      '+yhhhhy/+hhhhhs' -shhhyo-'      '+yhhy' '/yhhhhhhhhhhhyo/.
 ./oyo          '''         '---' -yhhhhh.    '''             .--     .--------.'
                                  yhhhhy.
                                  ohhhh/
                                   -+ss.
`

func (c *CanbusClient) showAboutDialog() *gocui.View {
  maxX, maxY := c.g.Size()

  var logoWithMargin string
  var i int = -1

  if (maxX-98) / 2 < 0 {
    logoWithMargin = logo
  } else {
    i = (maxX-98) / 2
    leftMargin := strings.Repeat(" ", i)
    logoWithMargin = strings.Replace(logo, "\n", fmt.Sprintf("\n%s", leftMargin), -1)
  }

  if v, err := c.g.SetView("about_dialog", -1, -1, maxX, maxY); err != nil {

    if err != gocui.ErrorUnkView {
      panic(err)
    }

    fmt.Fprintf(v, logoWithMargin)

    if err := c.g.SetCurrentView("about_dialog"); err != nil {
      panic(err)
    }

    return v
  }

  return nil
}
