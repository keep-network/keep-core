function terraform_fmt_check() {
  if [ $(terraform fmt -diff=true | grep -v .terraform | tee fmt_result.txt | wc -l) -gt 0 ]; then
    echo "Terraform formatting is not being used.:"
    echo
    cat fmt_result.txt
    rm fmt_result.txt
    git checkout -- .
    echo
    echo "Please run terraform fmt"
    exit 1
  fi
}

terraform_fmt_check