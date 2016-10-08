$containers = docker ps -a -q -f ancestor=blog:latest;

Write-Host ''

If ($containers -eq $null)
{
    Write-Host -ForegroundColor Gray 'There is nothing to remove. You are done here; go have fun!' 
}
Else {
    Write-Host -ForegroundColor Gray 'Found' $containers.Count 'container(s) to remove.'
    Write-Host -ForegroundColor Gray 'Sit back and relax - this may take a moment.'
    Write-Host ''
}

ForEach ($container in $containers)
{
    Write-Host 'Stopping' $container;
    Write-Host -NoNewline 'Stopped' $(docker stop $container);
    Write-Host -ForegroundColor Green ' ['$([char]8730)']'

    Write-Host 'Destroying' $container;
    Write-Host -NoNewLine 'Destroyed' $(docker rm $container);
    Write-Host -ForegroundColor Green ' ['$([char]8730)']'

    Write-Host -ForegroundColor Green 'Destruction of' $container 'should be complete.';
    Write-Host ''
}

Write-Host ''