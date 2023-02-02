module example

go 1.13

require gee v0.0.0
//在 go.mod 中使用 replace 将 gee 指向 ./gee
//引用相对路径
replace gee => ./gee
