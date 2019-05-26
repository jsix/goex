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
public class {{.table.Title}}Service {
    @Resource
    private {{$tableTitle}}Repository repo;

    public {{$tableTitle}}Entity save{{$tableTitle}}({{$tableTitle}}Entity {{.table.Name}}){
        return this.repo.save({{.table.Name}})
    }
    public void delete{{$tableTitle}}(int {{.pk}}) {
        this.repo.delete({{.pk}})
    }
    {{range $i,$c := .columns}}
    public {{$tableTitle}}Entity getBy{{$c.Title}}({{type "java" $c.TypeId}} {{$c.Name}}){
        return this.repo.findBy{{$c.Title}}({{$c.Name}})
    }
    {{end}}
}
