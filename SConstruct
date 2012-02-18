import os

gopath = os.path.abspath('.')
env = Environment(tools = ['default', 'gotool'])
env.PrependENVPath('GOPATH', gopath)
Export('env')

textproc = env.GoInstallPkg('textproc')
exe = env.GoInstall('pdfapp', textproc)
