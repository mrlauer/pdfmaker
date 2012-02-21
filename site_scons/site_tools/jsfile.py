from SCons.Script import *

def generate(env):
    coffee = Builder(
            action = "coffee -p -c $SOURCE > $TARGET",
            suffix = ".js",
            src_suffix = ".coffee",
            single_source = 1)
    env.Append(BUILDERS = { 'JSFile' : coffee })

def exists(env):
    return env.detect('JSFile')
