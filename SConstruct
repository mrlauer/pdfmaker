import os

gopath = os.path.abspath('gopath')
env = Environment(tools = ['default', 'gotool', 'jsfile'])
env.PrependENVPath('GOPATH', gopath)
env.PrependENVPath('PATH', ['/usr/local/bin', '/opt/node/bin'])
bindir = Dir('bin').abspath
env.SetDefault(BINDIR = bindir)

Export('env')

SConscript('gopath/SConscript')

# static files--javascript and whatnot
env.JSFile('static/scripts/main.js', 'staticsrc/scripts/main.coffee')
