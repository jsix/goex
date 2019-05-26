/**
 * this file is auto generated by GOF-CodeGenerator v{{.version}} !
 * if you want to modify this code,plese read guide doc
 * and modify code template later.
 *
 * guide please see http://s.to2.net/code-generator
 *
 */
package {{dot_pkg .global.Pkg}}.service

import {{dot_pkg .global.Pkg}}.pojo.{{.table.Title}}Entity
import {{dot_pkg .global.Pkg}}.repo.{{.table.Title}}Repository
import org.springframework.stereotype.Service
import javax.annotation.Resource
{{$tableTitle := .table.Title}}
/** {{.table.Comment}}服务  */
@Service
class {{.table.Title}}Service {

    @Resource
    lateinit var repo: {{$tableTitle}}Repository

    fun save{{$tableTitle}}({{.table.Name}}: {{$tableTitle}}Entity): {{$tableTitle}}Entity {
        return this.repo.save({{.table.Name}})
    }
    fun delete{{$tableTitle}}({{.pk}}:Int) {
         this.repo.delete({{.pk}})
    }
    {{range $i,$c := .columns}}
    fun getBy{{$c.Title}}({{$c.Name}}: {{type "kotlin" $c.TypeId}}):{{$tableTitle}}Entity? {
        return this.repo.findBy{{$c.Title}}({{$c.Name}})
    }
    {{end}}
}
