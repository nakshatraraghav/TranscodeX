import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as iam from "aws-cdk-lib/aws-iam";
import { env } from "../config/zenv";

export class EC2Instance extends Construct {
  public readonly instance: ec2.Instance;
  public readonly vpc: ec2.IVpc;
  public readonly sg: ec2.SecurityGroup;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    // Lookup the default VPC
    this.vpc = ec2.Vpc.fromLookup(this, "default_vpc", {
      vpcId: env.VPC_ID,
    });

    const keypair = ec2.KeyPair.fromKeyPairName(
      this,
      "transcodex_backend",
      "new-ssh-creds"
    );

    this.sg = new ec2.SecurityGroup(this, "SecurityGroup", {
      vpc: this.vpc,
      securityGroupName: "transcodex-backend-sg",
      allowAllOutbound: true,
      description: "Security group for allowing the backend server to communicate",
    });

    this.sg.addIngressRule(
      ec2.Peer.anyIpv4(), // Allow SSH access from any IP address
      ec2.Port.tcp(22),
      "Allow SSH access from anywhere"
    );

    this.sg.addIngressRule(
      ec2.Peer.anyIpv4(),
      ec2.Port.tcp(80),
      "Allow HTTP access from anywhere"
    );

    this.sg.addIngressRule(
      ec2.Peer.anyIpv4(),
      ec2.Port.tcp(443),
      "Allow HTTPS access from anywhere"
    );

    const role = new iam.Role(this, "EC2InstanceRole", {
      assumedBy: new iam.ServicePrincipal("ec2.amazonaws.com"),
    });

    role.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonS3FullAccess")
    );

    role.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonRDSFullAccess")
    );

    role.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonSQSFullAccess")
    );

    // Specify the subnet type as PUBLIC to ensure it gets a public IP
    this.instance = new ec2.Instance(this, "transcodex-backend-server-id", {
      instanceName: "transcodex-backend",
      machineImage: ec2.MachineImage.latestAmazonLinux2023({
        cpuType: ec2.AmazonLinuxCpuType.X86_64,
      }),
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T2,
        ec2.InstanceSize.MICRO
      ),
      keyName: keypair.keyPairName,
      vpc: this.vpc,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PUBLIC, // Ensure it's in a public subnet
      },
      securityGroup: this.sg,
      role: role,
      associatePublicIpAddress: true, // Ensure the instance has a public IP
    });
  }
}
