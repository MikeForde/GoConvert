# HL7 2.8 to MongoDB JSON Converter

A graphical application to convert **HL7 v2.8** messages into **MongoDB JSON** format compatible with IPS MERN/SERN applications. This tool allows users to easily paste or open HL7 message files and generate structured JSON outputs suitable for MongoDB storage.

## Features

- Converts **HL7 v2.8** messages into **MongoDB JSON** format.
- Graphical interface using the **Fyne** framework.
- Automatically suggests a filename based on the patient's name and package UUID.
- Supports saving the converted JSON file with a `.json` extension.
- Handles:
  - Patient information
  - Medications
  - Allergies
  - Conditions
  - Observations
  - Immunizations

## Installation

### Prerequisites

- **Go 1.18 or later** installed on your system.
- Fyne CLI for GUI packaging (optional for building standalone executables).

### Clone the Repository

```bash
git clone <repository-url>
cd <repository-directory>
