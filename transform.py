from neo4j import GraphDatabase
import pandas as pd
URI = "neo4j+s://2b7533a1.databases.neo4j.io"
AUTH = ("neo4j", "8V3hUCU44IYZwV3yMnZ-Vj4QgaLyg1N1pKlg3_b9UeE")
with GraphDatabase.driver(URI, auth=AUTH) as driver:
    driver.verify_connectivity()

NEO4J_URI="neo4j+s://2b7533a1.databases.neo4j.io"
NEO4J_USERNAME="neo4j"
NEO4J_PASSWORD="8V3hUCU44IYZwV3yMnZ-Vj4QgaLyg1N1pKlg3_b9UeE"
AURA_INSTANCEID="2b7533a1"
AURA_INSTANCENAME="Instance01"

# CSV file path
CSV_FILE = "./SemMed.csv"  # Replace with your CSV file path

class Neo4jGraph:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))
    
    def close(self):
        self.driver.close()
    
    def create_graph(self, data):
        with self.driver.session() as session:
            for _, row in data.iterrows():
                # Dynamically construct the query with the relationship type
                query = f"""
                MERGE (n1:`{row['Node1_Type']}` {{name: $node1}})
                MERGE (n2:`{row['Node2_Type']}` {{name: $node2}})
                MERGE (n1)-[r:{row['Relationship']} {{pmid: $pmid, description: $description}}]->(n2)
                """
                session.run(
                    query,
                    {
                        "node1": row["Node1"],
                        "node1_type": row["Node1_Type"],
                        "node2": row["Node2"],
                        "node2_type": row["Node2_Type"],
                        "pmid": row["PMID"],
                        "description": row["Description"]
                    }
                )

# Load the CSV data into a pandas DataFrame
data = pd.read_csv(CSV_FILE)

# Initialize the Neo4j connection
graph = Neo4jGraph(NEO4J_URI, NEO4J_USERNAME, NEO4J_PASSWORD)

# Create the graph
try:
    graph.create_graph(data)
    print("Graph created successfully!")
except Exception as e:
    print(f"An error occurred: {e}")
finally:
    graph.close()
