

local function change_context(ctx,name,code)
    ctx.name = name .. 'rewrite'
    ctx.code = code .. 'rewrite'
    return ctx
end