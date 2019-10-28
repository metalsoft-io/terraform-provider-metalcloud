#!/bin/bash

eval 'set +o history' 2>/dev/null || setopt HIST_IGNORE_SPACE 2>/dev/null
 touch ~/.gitcookies
 chmod 0600 ~/.gitcookies

 git config --global http.cookiefile ~/.gitcookies

 tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,FALSE,/,TRUE,2147483647,o,git-alexandru.bordei.gmail.com=1//03V1Ar-iXMplxCgYIARAAGAMSNwF-L9IrGLDmVa7NHZRP3sh4TGFz-5e8mnU3VjfK8QZspuBk4yBwSPbF0S2uktBuiKv3wPXaQC8
go-review.googlesource.com,FALSE,/,TRUE,2147483647,o,git-alexandru.bordei.gmail.com=1//03V1Ar-iXMplxCgYIARAAGAMSNwF-L9IrGLDmVa7NHZRP3sh4TGFz-5e8mnU3VjfK8QZspuBk4yBwSPbF0S2uktBuiKv3wPXaQC8
__END__
eval 'set -o history' 2>/dev/null || unsetopt HIST_IGNORE_SPACE 2>/dev/null


