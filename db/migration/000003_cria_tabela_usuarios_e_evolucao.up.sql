CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.usu_usuario (
    usu_tx_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usu_tx_nome VARCHAR(100) NOT NULL,
    usu_tx_email VARCHAR(100) UNIQUE NOT NULL,
    usu_tx_senha VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS treino.evo_evolucao (
    evo_nr_id SERIAL PRIMARY KEY,
    usu_tx_id UUID NOT NULL,
    evo_dt_data DATE DEFAULT CURRENT_DATE,
    evo_nr_peso DECIMAL(5,2),
    evo_nr_altura DECIMAL(4,2),
    evo_nr_ombro DECIMAL(5,2),
    evo_nr_busto DECIMAL(5,2),
    evo_nr_abdomen DECIMAL(5,2),
    evo_nr_cintura DECIMAL(5,2),
    evo_nr_quadril DECIMAL(5,2),
    evo_nr_braco_direito DECIMAL(5,2),
    evo_nr_braco_esquerdo DECIMAL(5,2),
    evo_nr_antebraco_direito DECIMAL(5,2),
    evo_nr_antebraco_esquerdo DECIMAL(5,2),
    evo_nr_coxa_direita DECIMAL(5,2),
    evo_nr_coxa_esquerda DECIMAL(5,2),
    evo_nr_panturrilha_direita DECIMAL(5,2),
    evo_nr_panturrilha_esquerda DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_evolucao_usuario FOREIGN KEY (usu_tx_id) REFERENCES auth.usu_usuario(usu_tx_id)
);

ALTER TABLE treino.tre_treino
ADD COLUMN IF NOT EXISTS usu_tx_id UUID,
ADD CONSTRAINT fk_treino_usuario FOREIGN KEY (usu_tx_id) REFERENCES auth.usu_usuario(usu_tx_id);
