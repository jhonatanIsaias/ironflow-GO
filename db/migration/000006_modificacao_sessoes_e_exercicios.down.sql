-- Remove coluna exe_tm_hora_fim de exe_exercicio
ALTER TABLE treino.exe_exercicio 
DROP COLUMN exe_tm_hora_fim;

-- Remove coluna set_dt_data_fim de set_sessao_treino
ALTER TABLE treino.set_sessao_treino 
DROP COLUMN set_dt_data_fim;

-- Remove coluna set_tm_hora_fim de set_sessao_treino
ALTER TABLE treino.set_sessao_treino 
DROP COLUMN set_tm_hora_fim;

-- Remove coluna exe_tx_tipo_equipamento de fit_ficha_treino
ALTER TABLE treino.fit_ficha_treino 
DROP COLUMN exe_tx_tipo_equipamento;

-- Adiciona coluna exe_tx_tipo_equipamento em exe_exercicio
ALTER TABLE treino.exe_exercicio 
ADD COLUMN exe_tx_tipo_equipamento VARCHAR(50);
