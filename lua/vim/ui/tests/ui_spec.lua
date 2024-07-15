local describe = require("busted").describe
local it = require("busted").it
local before_each = require("busted").before_each
local assert = require("busted").assert

local eq = assert.are.same
local float = require("lua.vim.ui.float")
local utils = require("lua.vim.tests.utils")

describe("harpoon", function ()
  before_each(function ()
    utils.clear_memory()
    float = require("lua.vim.ui.float")
  end)

  it("window toggle test", function ()
    eq(float.buf_id, nil)
    eq(float.win_id, nil)

    float:toggle()

    local win_id = float.win_id
    local buf_id = float.buf_id
    eq(true, vim.api.nvim_win_is_valid(win_id))
    eq(true, vim.api.nvim_buf_is_valid(float.buf_id))

    float:toggle()

    eq(float.buf_id, nil)
    eq(float.win_id, nil)
    eq(false, vim.api.nvim_win_is_valid(win_id))
    eq(false, vim.api.nvim_buf_is_valid(float.buf_id))
  end)

  it("window toggle test", function ()
    eq(float.buf_id, nil)
    eq(float.win_id, nil)

    float:toggle()

    local win_id = float.win_id
    local buf_id = float.buf_id
    eq(true, vim.api.nvim_win_is_valid(win_id))
    eq(true, vim.api.nvim_buf_is_valid(float.buf_id))

    vim.api.nvim_buf_delete(buf_id, { force = true })

    eq(nil, float.buf_id)
    eq(nil, float.win_id)
    eq(false, vim.api.nvim_win_is_valid(win_id))
    eq(false, vim.api.nvim_buf_is_valid(float.buf_id))
  end)
end)
