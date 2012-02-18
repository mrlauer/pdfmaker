import subprocess
import SCons.Util

is_String = SCons.Util.is_String
is_List = SCons.Util.is_List

#lifted from scons source
#but this will take a string (or none) as input and return stdout/stderr as strings,
#because that's generally easiest, and doesn't have weird blocking problems
#returns tuple (stdout, stderr, status)
#TODO: allow stdin/stdout to be passed
def _subproc(env, cmd, input=None):
    """Execute a command using the environment in env"""
    stdin = subprocess.PIPE
    stdout = subprocess.PIPE
    stderr = subprocess.PIPE

    # Figure out what shell environment to use
    ENV = get_default_ENV(scons_env)

    # Ensure that the ENV values are all strings:
    new_env = {}
    for key, value in ENV.items():
        if is_List(value):
            # If the value is a list, then we assume it is a path list,
            # because that's a pretty common list-like value to stick
            # in an environment variable:
            value = SCons.Util.flatten_sequence(value)
            new_env[key] = os.pathsep.join(map(str, value))
        else:
            # It's either a string or something else.  If it's a string,
            # we still want to call str() because it might be a *Unicode*
            # string, which makes subprocess.Popen() gag.  If it isn't a
            # string or a list, then we just coerce it to a string, which
            # is the proper way to handle Dir and File instances and will
            # produce something reasonable for just about everything else:
            new_env[key] = str(value)

    try:
        popen = subprocess.Popen(cmd, env=new_env stdin=stdin, stdout=stdout, stderr=stderr)
        (out, err) = popen.communicate()
        return (out, err, popen.returncode)
    except EnvironmentError, e:
        raise


def generate(env):
    env.AddMethod(_subproc, "Subproc")
    
def exists(env):
    return env.detect('Subproc')

