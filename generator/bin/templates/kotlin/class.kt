/**
 * this file is auto generated by GOF-CodeGenerator v{{.version}} !
 * if you want to modify this code,plese read guide doc
 * and modify code template later.
 *
 * guide please see http://s.to2.net/code-generator
 *
 */

package {{dot_pkg .global.Pkg}};

/** {{.T.Comment}} */
class {{.T.Title}}{
    {{range $i,$c:=.columns}}
    /** {{$c.Comment}} */
    var {{lowerTitle $c.Title}}:{{type "kotlin" $c.TypeId}} = {{end}}
}