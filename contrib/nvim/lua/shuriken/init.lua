-- shuriken in neovim: what the scrolls hit, shown as diagnostics beside
-- gopls' own. On opening or saving a Go file, shuriken throws at the
-- whole module, and every loaded Go buffer in that module gets the hits
-- in its file: signs in the gutter, the words at the end of the line,
-- and ]d, [d and the line float like any other diagnostic.

local M = {}

-- namespace keeps shuriken's diagnostics apart from gopls', so neither
-- clears the other's.
local namespace = vim.api.nvim_create_namespace('shuriken')

-- options are what setup may change.
local options = {
  -- cmd is the shuriken to run; the one on PATH unless given.
  cmd = 'shuriken',
  -- severity is how every hit is shown.
  severity = vim.diagnostic.severity.WARN,
}

-- throwing holds the module roots with a throw in flight, so a save
-- during a throw does not start a second one.
local throwing = {}

-- moduleRoot is the directory of the go.mod above path, or nil.
local function moduleRoot(path)
  local found = vim.fs.find('go.mod', { upward = true, path = vim.fs.dirname(path) })[1]
  return found and vim.fs.dirname(found) or nil
end

-- diagnosticsByFile turns shuriken's facts into diagnostics, keyed by
-- each file's full path.
local function diagnosticsByFile(root, facts)
  local byFile = {}
  for _, hit in ipairs(facts.hits or {}) do
    local file = vim.fs.normalize(root .. '/' .. hit.file)
    byFile[file] = byFile[file] or {}
    table.insert(byFile[file], {
      lnum = math.max((hit.line or 1) - 1, 0),
      col = math.max((hit.column or 1) - 1, 0),
      severity = options.severity,
      source = 'shuriken',
      code = hit.scroll,
      message = hit.scroll .. ': ' .. hit.seen,
    })
  end
  return byFile
end

-- show sets each loaded Go buffer under root to the hits in its file,
-- and clears the ones with none.
local function show(root, byFile)
  for _, buffer in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_loaded(buffer) then
      local name = vim.fs.normalize(vim.api.nvim_buf_get_name(buffer))
      if name:sub(-3) == '.go' and moduleRoot(name) == root then
        vim.diagnostic.set(namespace, buffer, byFile[name] or {})
      end
    end
  end
end

-- answered reads what shuriken said. shuriken exits 1 when a scroll is
-- broken and still writes its facts, so the json decides, not the code.
local function answered(root, done)
  local ok, facts = pcall(vim.json.decode, done.stdout or '')
  if ok and type(facts) == 'table' then
    show(root, diagnosticsByFile(root, facts))
    return
  end
  local said = vim.trim(done.stderr or '')
  vim.notify('shuriken: ' .. (said ~= '' and said or 'no facts came back'),
    vim.log.levels.WARN)
end

-- M.throw runs shuriken at the module that buffer's file is in.
function M.throw(buffer)
  local path = vim.api.nvim_buf_get_name(buffer)
  local root = path ~= '' and moduleRoot(path) or nil
  if not root or throwing[root] then
    return
  end
  throwing[root] = true
  local started = pcall(vim.system, { options.cmd, 'facts', '.' },
    { cwd = root, text = true }, function(done)
      vim.schedule(function()
        throwing[root] = nil
        answered(root, done)
      end)
    end)
  if not started then
    throwing[root] = nil
    vim.notify('shuriken: ' .. options.cmd .. ' could not be run',
      vim.log.levels.WARN)
  end
end

-- M.setup sets the options given and throws on opening and saving Go
-- files.
function M.setup(given)
  options = vim.tbl_extend('force', options, given or {})
  local group = vim.api.nvim_create_augroup('shuriken', { clear = true })
  vim.api.nvim_create_autocmd({ 'BufReadPost', 'BufWritePost' }, {
    group = group,
    pattern = '*.go',
    callback = function(event) M.throw(event.buf) end,
  })
end

return M
