import os.path
from SCons.Script import *
import SCons.Util

# Install javascript files
def InstallJS(env, targetdir, source, sourcedir = None):
    results = []
    if source is not None and not SCons.Util.is_List(source):
        source = [source]
    if not sourcedir:
        sourcedir = Dir('.').abspath
    for f in source:
        rel = os.path.relpath(str(f), sourcedir)
        target = os.path.join(targetdir, rel)
        base, ext = os.path.splitext(target)
        target = base + '.js'
        results += env._InstallJS(target, f)

def generate(env):
    coffee = Builder(
            action = "coffee -p -c $SOURCE > $TARGET",
            suffix = ".js",
            src_suffix = ".coffee",
            single_source = 1)
    env.Append(BUILDERS = { 'JSFile' : coffee })

    # install one file, however it works
    installer = Builder(action = {}, suffix = '.js')
    installer.add_action('.js', action = Copy('$TARGET', '$SOURCE'))
    installer.add_action('.coffee', 
            action = 'coffee -p -c $SOURCE > $TARGET')
    env.Append(BUILDERS = { '_InstallJS' : installer })
    env.AddMethod(InstallJS)

def exists(env):
    return env.detect('JSFile')
