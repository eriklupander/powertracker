import {Construct, RemovalPolicy} from "@aws-cdk/core";
import {CfnDatabase, CfnTable} from "@aws-cdk/aws-timestream";

export interface imestreamConstructProps {
    databaseName: string;
    tableName: string;
}

export class TimestreamConstruct extends Construct {
    public readonly database: CfnDatabase;
    public readonly table: CfnTable;

    constructor(scope: Construct, id: string, props: imestreamConstructProps) {
        super(scope, id);
        this.database = new CfnDatabase(this, 'Database', {
            databaseName: props.databaseName,
        });
        this.database.applyRemovalPolicy(RemovalPolicy.RETAIN);

        this.table = new CfnTable(this, 'Table', {
            tableName: props.tableName,
            databaseName: props.databaseName,
            retentionProperties: {
                memoryStoreRetentionPeriodInHours: (48).toString(10), // 2 days
                magneticStoreRetentionPeriodInDays: (365 * 2).toString(10) // 2 years
            }
        });
        this.table.node.addDependency(this.database);
        this.table.applyRemovalPolicy(RemovalPolicy.RETAIN);
    }
}