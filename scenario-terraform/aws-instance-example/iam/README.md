## Terraform IAM AWS


### Checklist

1. Region must AWS-EC2 ap-southeast-1
2. Stick to the naming module rules, (group, user, etc)
3. Check password_reset_required == true
4. Check attach_iam_self_management_policy = false
5. Check if user access to all resource EC2