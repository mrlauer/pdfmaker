import os

extGopath = os.environ.get('GOPATH')
gopath = os.path.abspath('gopath')
env = Environment(tools = ['default', 'gotool', 'jsfile'])
if extGopath:
    env.PrependENVPath('GOPATH', extGopath)
env.PrependENVPath('GOPATH', gopath)
env.PrependENVPath('PATH', ['/usr/local/bin', '/opt/node/bin'])
bindir = Dir('bin').abspath
env.SetDefault(BINDIR = bindir)
env.SetDefault(STATICDIR = Dir('static').abspath)

Export('env')

SConscript('gopath/SConscript')

# static files--javascript and whatnot
SConscript('staticsrc/SConscript')

run = env.Alias('run', ['bin', 'static'], './bin/pdfapp')
AlwaysBuild(run)
