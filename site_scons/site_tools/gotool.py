import os.path
import platform
import glob
from SCons.Script import *

def _go_generator(source, target, env, for_signature):
    """
    Generator for go actions.
    Get the name from the basename of the target.
    TODO: this needs to handle nested paths.
    """
    path = os.path.basename(str(target[0]))
    t, _ = os.path.splitext(path)
    return "$GOINSTALL %s" % t

def GoInstall(env, name, deps = None):
    # get the gopath
    gopath = env['ENV']['GOPATH']
    # TODO: handle lists
    sourcedir = os.path.join(gopath, 'src', name)
    files = glob.glob(os.path.join(sourcedir, '*.go'))
    files += glob.glob(os.path.join(sourcedir, '*.c'))
    if deps:
        files += deps
    target = os.path.join(gopath, 'bin', name + env.subst('$PROGSUFFIX'))
    return env._GoInstall(target = target, source = files)

def GoInstallPkg(env, name, deps = None):
    # get the gopath
    gopath = env['ENV']['GOPATH']
    # TODO: handle lists
    sourcedir = os.path.join(gopath, 'src', name)
    files = glob.glob(os.path.join(sourcedir, '*.go'))
    files += glob.glob(os.path.join(sourcedir, '*.c'))
    if deps:
        files += deps
    pkgdir = env.subst('${GOOS}_${GOARCH}')
    target = os.path.join(gopath, 'pkg', 
            pkgdir, name + env.subst('$LIBSUFFIX'))
    return env._GoInstall(target = target, source = files)

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
    
    goinstall = Builder(
            generator = _go_generator)
    env.Append(BUILDERS = { '_GoInstall' : goinstall})
    env.AddMethod(GoInstall)
    env.AddMethod(GoInstallPkg)

def exists(env):
    return env.detect('GoInstall')
