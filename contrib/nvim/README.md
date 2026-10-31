# shuriken in neovim

shuriken's hits, shown as diagnostics beside gopls' own: a sign in the
gutter, the words at the end of the line, `]d` and `[d` to move between
them, and the line float to read one whole. Each shows its scroll by name
and says it came from shuriken.

On opening or saving a Go file, shuriken throws at the whole module (it
always does), and every open Go buffer in that module is brought up to
date. A save during a throw waits for the throw in flight.

With lazy.nvim, one spec in `lua/plugins/`:

    return {
      dir = '~/path/to/libninja/contrib/nvim',
      ft = 'go',
      config = function() require('shuriken').setup() end,
    }

`setup` takes `cmd`, the shuriken to run (the one on PATH unless given),
and `severity`, a `vim.diagnostic.severity` (warnings unless given).

It needs neovim 0.10 or later, for `vim.system`.
