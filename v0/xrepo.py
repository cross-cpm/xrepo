#!/usr/bin/env python


# reference from https://github.com/develersrl/git-externals.git


import sys
import subprocess
import os
import re
import os.path
import posixpath
import json

from subprocess import check_call
from contextlib import contextmanager

try:
    from urllib.parse import urlparse, urlsplit, urlunsplit
except ImportError:
    from urlparse import urlparse, urlsplit, urlunsplit





EXTERNALS_JSON = 'externals.json'
EXTERNALS_ROOT = '.externals'



################################################################################


def echo(*args):
    print(args)

def info(*args):
    print(args)

def error(*args):
    print(args)


################################################################################


class ProgError(Exception):
    def __init__(self, prog='', errcode=1, errmsg='', args=''):
        if isinstance(args, tuple):
            args = u' '.join(args)
        super(ProgError, self).__init__(u'\"{} {}\" {}'.format(prog, args, errmsg))
        self.prog = prog
        self.errcode = errcode

    def __str__(self):
        name = u'{}Error'.format(self.prog.title())
        msg = super(ProgError, self).__str__()
        return u'<{}: {} {}>'.format(name, self.errcode, msg)


class GitError(ProgError):
    def __init__(self, **kwargs):
        super(GitError, self).__init__(prog='git', **kwargs)


class SvnError(ProgError):
    def __init__(self, **kwargs):
        super(SvnError, self).__init__(prog='svn', **kwargs)


class GitSvnError(ProgError):
    def __init__(self, **kwargs):
        super(GitSvnError, self).__init__(prog='git-svn', **kwargs)


class CommandError(ProgError):
    def __init__(self, cmd, **kwargs):
        super(CommandError, self).__init__(prog=cmd, **kwargs)


def svn(*args, **kwargs):
    universal_newlines = kwargs.get('universal_newlines', True)
    output, err, errcode = _command('svn', *args, capture=True, universal_newlines=universal_newlines)
    if errcode != 0:
        print("running svn ", args)
        raise SvnError(errcode=errcode, errmsg=err)
    return output


def git(*args, **kwargs):
    capture = kwargs.get('capture', True)
    output, err, errcode = _command('git', *args, capture=capture, universal_newlines=True)
    if errcode != 0:
        raise GitError(errcode=errcode, errmsg=err, args=args)
    return output


def gitsvn(*args, **kwargs):
    capture = kwargs.get('capture', True)
    output, err, errcode = _command('git', 'svn', *args, capture=capture, universal_newlines=True)
    if errcode != 0:
        raise GitSvnError(errcode=errcode, errmsg=err, args=args)
    return output


def gitsvnrebase(*args, **kwargs):
    capture = kwargs.get('capture', True)
    output, err, errcode = _command('git-svn-rebase', *args, capture=capture, universal_newlines=True)
    if errcode != 0:
        raise GitSvnError(errcode=errcode, errmsg=err, args=args)
    return output


def command(cmd, *args, **kwargs):
    universal_newlines = kwargs.get('universal_newlines', True)
    capture = kwargs.get('capture', True)
    output, err, errcode = _command(cmd, *args, universal_newlines=universal_newlines, capture=capture)
    if errcode != 0:
        raise CommandError(cmd, errcode=errcode, errmsg=err, args=args)
    return output


def _command(cmd, *args, **kwargs):
    env = kwargs.get('env', dict(os.environ))
    env.setdefault('LC_MESSAGES', 'C')
    universal_newlines = kwargs.get('universal_newlines', True)
    capture = kwargs.get('capture', True)
    if capture:
        stdout, stderr = subprocess.PIPE, subprocess.PIPE
    else:
        stdout, stderr = None, None

    p = subprocess.Popen([cmd] + list(args),
                         stdout=stdout,
                         stderr=stderr,
                         universal_newlines=universal_newlines,
                         env=env)
    output, err = p.communicate()
    return output, err, p.returncode


def current_branch():
    return git('name-rev', '--name-only', 'HEAD').strip()


def branches():
    refs = git('for-each-ref', 'refs/heads', "--format=%(refname)")
    return [line.split('/')[2] for line in refs.splitlines()]


def tags():
    refs = git('for-each-ref', 'refs/tags', "--format=%(refname)")
    return [line.split('/')[2] for line in refs.splitlines()]


TAGS_RE = re.compile('.+/tags/(.+)')

def git_remote_branches_and_tags():
    output = git('branch', '-r')

    _branches, _tags = [], []

    for line in output.splitlines():
        line = line.strip()
        m = TAGS_RE.match(line)

        t = _tags if m is not None else _branches
        t.append(line)

    return _branches, _tags


@contextmanager
def checkout(branch, remote=None, back_to='master', force=False):
    brs = set(branches())

    cmd = ['git', 'checkout']
    if force:
        cmd += ['--force']
    # if remote is not None -> create local branch from remote
    if remote is not None and branch not in brs:
        check_call(cmd + ['-b', branch, remote])
    else:
        check_call(cmd + [branch])
    yield
    check_call(cmd + [back_to])


@contextmanager
def chdir(path):
    cwd = os.path.abspath(os.getcwd())

    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(cwd)


def mkdir_p(path):
    if path != '' and not os.path.exists(path):
        os.makedirs(path)


def header(msg):
    banner = '=' * 78

    print('')
    print(banner)
    print(u'{:^78}'.format(msg))
    print(banner)


def print_msg(msg):
    print(u'  {}'.format(msg))



################################################################################



def externals_json_path(pwd=None):
    return os.path.join(pwd or root_path(), EXTERNALS_JSON)


def root_path():
    return git('rev-parse', '--show-toplevel').strip()


def externals_root_path(pwd=None):
    return os.path.join(pwd or root_path(), EXTERNALS_ROOT)


def get_repo_name(repo):
    externals = load_gitexts()
    if repo in externals and 'name' in externals[repo]:
        # echo ("for {} in pwd:{} returning {}".format(repo, os.getcwd(),
        #                                              externals[repo]['name']))
        return externals[repo]['name']

    if repo[-1] == '/':
        repo = repo[:-1]
    name = repo.split('/')[-1]
    if name.endswith('.git'):
        name = name[:-len('.git')]
    if not name:
        error("Invalid repository name: \"{}\"".format(repo), exitcode=1)
    return name


def load_gitexts(pwd=None):
    """Load the *externals definition file* present in given
    directory, or cwd
    """
    d = pwd if pwd is not None else '.'
    fn = os.path.join(d, EXTERNALS_JSON)
    if os.path.exists(fn):
        with open(fn) as f:
            return normalize_gitexts(json.load(f))
    return {}


def normalize_gitexts(gitext):
    for url, repo_data in gitext.items():
        # svn external url must be absolute and svn+ssh to be autodetected
        gitext[url].setdefault('vcs', 'svn' if 'svn' in urlparse(url).scheme else 'git')

        for (src, dsts) in repo_data['targets'].items():
            if src == './':
                for dst in dsts:
                    workdir = os.path.join(os.getcwd(), dst.replace('/', os.path.sep))
                    gitext[url]['workdir'] = workdir
                    gitext[url]['name'] = os.path.basename(workdir)
                    break

    return gitext


def normalize_gitext_url(url):
    # an absolute url is already normalized
    if urlparse(url).netloc != '' or url.startswith('git@'):
        return url

    # relative urls use the root url of the current origin
    remote_name = git('config', 'branch.%s.remote' % current_branch()).strip()
    remote_url = git('config', 'remote.%s.url' % remote_name).strip()

    if remote_url.startswith('git@'):
        prefix = remote_url[:remote_url.index(':')+1]
        remote_url = prefix + url.strip('/')
    else:
        remote_url = urlunsplit(urlsplit(remote_url)._replace(path=url))

    return remote_url


def filter_externals_not_needed(all_externals, entries):
    git_externals = {}
    for repo_name, repo_val in all_externals.items():
        filtered_targets = {}
        for src, dsts in repo_val['targets'].items():
            filtered_dsts = []
            for dst in dsts:
                inside_external = any([os.path.abspath(dst).startswith(e) for e in entries])
                if inside_external:
                    filtered_dsts.append(dst)

            if filtered_dsts:
                filtered_targets[src] = filtered_dsts

        if filtered_targets:
            git_externals[repo_name] = all_externals[repo_name]
            git_externals[repo_name]['targets'] = filtered_targets

    return git_externals


def sparse_checkout(repo_name, repo, dirs):
    git('init', repo_name)

    with chdir(repo_name):
        git('remote', 'add', '-f', 'origin', repo)
        git('config', 'core.sparsecheckout', 'true')

        with open(os.path.join('.git', 'info', 'sparse-checkout'), 'wt') as fp:
            fp.write('{}\n'.format(externals_json_path()))
            for d in dirs:
                # assume directories are terminated with /
                fp.write(posixpath.normpath(d))
                if d[-1] == '/':
                    fp.write('/')
                fp.write('\n')

    return repo_name


def resolve_revision(ref, mode='git'):
    assert mode in ('git', 'svn'), "mode = {} not in (git, svn)".format(mode)
    if ref is not None:
        if ref.startswith('svn:r'):
            # echo("Resolving {}".format(ref))
            ref = ref.strip('svn:r')
            # If the revision starts with 'svn:r' in 'git' mode we search
            # for the matching hash.
            if mode == 'git':
                ref = git('log', '--grep', 'git-svn-id:.*@%s' % ref, '--format=%H', capture=True).strip()
    return ref


def get_work_parent_dir(repo_data):
    return os.path.dirname(repo_data['workdir']) if repo_data.get('workdir') else externals_root_path()


def gitext_up(reset=False, use_gitsvn=False):

    if not os.path.exists(externals_json_path()):
        return

    git_externals = load_gitexts()

    def egit(command, *args):
        if command == 'checkout' and reset:
            args = ('--force',) + args
        git(command, *args, capture=False)

    def git_initial_checkout(repo_name, repo_url):
        """Perform the initial git clone (or sparse checkout)"""
        dirs = git_externals[ext_repo]['targets'].keys()
        if './' not in dirs:
            echo('Doing a sparse checkout of:', ', '.join(dirs))
            sparse_checkout(repo_name, repo_url, dirs)
        else:
            egit('clone', repo_url, repo_name)

    def git_update_checkout(reset):
        """Update an already existing git working tree"""
        if reset:
            egit('reset', '--hard')
            egit('clean', '-df')
        egit('fetch', '--all')
        egit('fetch', '--tags')
        if 'tag' in git_externals[ext_repo]:
            echo('Checking out tag', git_externals[ext_repo]['tag'])
            egit('checkout', git_externals[ext_repo]['tag'])
        else:
            echo('Checking out branch', git_externals[ext_repo]['branch'])
            egit('checkout', git_externals[ext_repo]['branch'])

            rev = get_rev(ext_repo)
            if rev is not None:
                echo('Checking out commit', rev)
                egit('checkout', rev)

    def get_rev(ext_repo, mode='git'):
        ref = git_externals[ext_repo]['ref']
        return resolve_revision(ref, mode)

    def gitsvn_initial_checkout(repo_name, repo_url):
        """Perform the initial git-svn clone (or sparse checkout)"""
        min_rev = get_rev(ext_repo, mode='svn') or 'HEAD'
        gitsvn('clone', normalized_ext_repo, repo_name, '-r%s' % min_rev, capture=False)

    def gitsvn_update_checkout(reset):
        """Update an already existing git-svn working tree"""
        # FIXME: seems this might be necessary sometimes (happened with
        # 'vectorfonts' for example that the following error: "Unable to
        # determine upstream SVN information from HEAD history" was fixed by
        # adding that, but breaks sometimes. (investigate)
        # git('rebase', '--onto', 'git-svn', '--root', 'master')
        gitsvnrebase('.', capture=False)
        rev = get_rev(ext_repo) or 'git-svn'
        echo('Checking out commit', rev)
        git('checkout', rev)

    def svn_initial_checkout(repo_name, repo_url):
        """Perform the initial svn checkout"""
        svn('checkout', '--ignore-externals', normalized_ext_repo, repo_name, capture=False)

    def svn_update_checkout(reset):
        """Update an already existing svn working tree"""
        if reset:
            svn('revert', '-R', '.')
        rev = get_rev(ext_repo, mode='svn') or 'HEAD'
        echo('Updating to commit', rev)
        svn('up', '--ignore-externals', '-r%s' % rev, capture=False)

    def autosvn_update_checkout(reset):
        if os.path.exists('.git'):
            gitsvn_update_checkout(reset)
        else:
            svn_update_checkout(reset)


    for ext_repo in git_externals.keys():
        normalized_ext_repo = normalize_gitext_url(ext_repo)

        if git_externals[ext_repo]['vcs'] == 'git':
            _initial_checkout = git_initial_checkout
            _update_checkout = git_update_checkout
        else:
            if use_gitsvn:
                _initial_checkout = gitsvn_initial_checkout
            else:
                _initial_checkout = svn_initial_checkout
            _update_checkout = autosvn_update_checkout

        work_parent_dir = get_work_parent_dir(git_externals[ext_repo])
        # print "work_parent_dir", work_parent_dir

        mkdir_p(work_parent_dir)
        with chdir(work_parent_dir):
            repo_name = get_repo_name(normalized_ext_repo)
            ext_name = git_externals[ext_repo].get('name', '')
            ext_name = ext_name if ext_name else repo_name

            info('External', ext_name)
            if not os.path.exists(ext_name):
                echo('Cloning external', ext_name)
                _initial_checkout(ext_name, normalized_ext_repo)

            with chdir(ext_name):
                echo('Retrieving changes from server: ', ext_name)
                _update_checkout(reset)



################################################################################


def dump_usage():
    echo(sys.argv[0], '<subcommand>')


def main(argv):
    if len(argv) < 1:
        dump_usage()
        return

    subcmd = argv[0]
    if subcmd == "up" or subcmd == "update":
        gitext_up()
    else:
        error("unsupport sub command", subcmd)


if __name__ == "__main__":
    main(sys.argv[1:])
