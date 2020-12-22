/* Project by Dylan Luttrell (2020).

With inspiration from phpSyntaxTree,
by André Eisenbach and Mei Eisenbach
(2003, http://www.tycho.iel.unicamp.br/phpsyntaxtree/)
*/
package main

import (
	"github.com/APDevice/syntax_tree/lib"
)

func main() {
	no, err := lib.NewSentence("[S [NP [N goSyntaxTree]][VP [V makes][NP [AdjP [Adj awesome]] [AdjP [Adj syntax]] [N trees]]]]")

	if err == nil {
		no.Render()

	}
}
