import os

gopath = os.path.abspath('.')
env = Environment(tools = ['default', 'gotool', 'jsfile'])
env.PrependENVPath('GOPATH', gopath)
env.PrependENVPath('PATH', ['/usr/local/bin', '/opt/node/bin'])
Export('env')

textproc = env.GoInstallPkg('textproc')
exe = env.GoInstall('pdfapp', textproc)

# static files--javascript and whatnot
env.JSFile('static/scripts/main.js', 'staticsrc/scripts/main.coffee')
