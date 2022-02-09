--============================================================
-- @Author: Mark(车英达)
-- @E-mail: 120006937@qq.com
-- @Date:   2019-10-10 14:30:08
-- @Last Modified time: 2019-10-25 16:11:29
--============================================================

local data = {
		name = "cyd",
		age = "28",
		nums = {1,2,3}
}
function hello(  )
	print("hello world")
end

function get_str( ... )
	local args = {...}
	local str = "hello "
	for i,v in ipairs(args) do
		str = str .. v
	end
	return str
end

function get_table(  )
	return data
end

function get_num(  )
	return 10.1
end

function test_call_go( d )
	d.Set("aaa", "30")
end