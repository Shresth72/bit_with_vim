local float = require("lua.vim.ui.float")

local M = {}

function M.clear_memory()
  if float.buf_id ~= nil then
    float:toggle()
  end

  -- require("plenary").reload.reload_module("vim_api")
end

return M
