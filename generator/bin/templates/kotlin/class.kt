/**
 * this file is auto generated by GOF-CodeGenerator v{{.version}} !
 * if you want to modify this code,please read guide doc
 * and modify code template later.
 *
 * guide please see http://s.to2.net/code-generator
 *
 */

package {{pkg "go" .global.Pkg}};

/** {{.table.Comment}} */
class {{.table.Title}}{
    {{range $i,$c:=.columns}}
    /** {{$c.Comment}} */
    var {{lowerTitle $c.Title}}:{{type "kotlin" $c.TypeId}} = {{end}}
}