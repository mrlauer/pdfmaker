import os.path
import platform
import glob
from SCons.Script import *

def _findSrc(gopath, name):
    for path in gopath.split(':'):
        candidate = os.path.join(path, 'src', name)
        if os.path.exists(candidate):
            return candidate, path
    return None, None

def _first(gopath):
    return gopath.split(':')[0]

def _goinstall(env, out, suffix, name, deps):
    # get the gopath
    gopath = env['ENV']['GOPATH']
    # TODO: handle lists
    sourcedir, srcgopath = _findSrc(gopath, name)
    files = glob.glob(os.path.join(sourcedir, '*.go'))
    files += glob.glob(os.path.join(sourcedir, '*.c'))
    if deps:
        files += deps
    target = os.path.join(_first(gopath), out, name + suffix)
    cmds = ["$GOINSTALL $GOINSTALLFLAGS %s" % name ]
    #cmds += [Move(target, os.path.join(srcgopath, out, name + suffix))]
    return env.Command(action = cmds, target = target, source = files)

def GoInstall(env, name, deps = None):
    return _goinstall(env, 'bin', env.subst('$PROGSUFFIX'), name, deps)

def GoInstallPkg(env, name, deps = None):
    pkgdir = env.subst('${GOOS}_${GOARCH}')
    pkg = os.path.join('pkg', pkgdir)
    return _goinstall(env, pkg, env.subst('$LIBSUFFIX'), name, deps)

def generate(env):
    # figure out the environment
    # Should this be in a Configure block?
    # we'll get GOROOT from the environment for now.
    env.SetDefault(GOROOT = os.environ.get('GOROOT'))
    system = platform.system()
    env.SetDefault(GOOS = system.lower())
    machine = platform.machine()
    goarch = "amd64"
    goarchchar = "6"
    is64 = (sys.maxsize > 2**32)
    if not is64:
        if machine == "i386":
            goarch = "386"
            goarchchar = "8"
        else:
            goarch = "arm"
            goarchchar = "5"
    env.SetDefault(GOARCH = goarch, GOARCHCHAR = goarchchar) 

    env.PrependENVPath('PATH', os.path.join(os.environ['GOROOT'], 'bin'))

    env.SetDefault(GOINSTALL = 'go install')
    
    env.AddMethod(GoInstall)
    env.AddMethod(GoInstallPkg)

def exists(env):
    return env.detect('GoInstall')
